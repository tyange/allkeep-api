package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tyange/white-shadow-api/models"
	"github.com/tyange/white-shadow-api/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

var (
	oauthConfig *oauth2.Config
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindBodyWithJSON(&user)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.Save()

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not save user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}

func login(context *gin.Context) {
	fmt.Println("hi")

	var user models.User

	err := context.ShouldBindBodyWithJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not authenticate user."})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	context.SetCookie("session", token, 60*60*2, "/", "localhost", false, false)
	context.JSON(http.StatusOK, gin.H{"message": "Login successful!"})
}

func googleLoginCallBack(context *gin.Context) {
	var code map[string]string
	err := context.ShouldBindBodyWithJSON(&code)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get a code from body."})
		return
	}

	googleLoginToken, err := oauthConfig.Exchange(context, code["code"])
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get a token from code."})
		return
	}

	client := oauthConfig.Client(context, googleLoginToken)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get a user info from token."})
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not reading user info body."})
		return
	}

	var userInfo UserInfo
	err = json.Unmarshal(body, &userInfo)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error unmarshaling user info to JSON."})
		return
	}

	var user models.User
	user.Email = userInfo.Email

	isDuplicated := user.CheckDuplicateUserId()

	if !isDuplicated {
		err := user.SaveWithoutPassword()

		if err != nil {
			fmt.Println(err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not created user with google login."})
			return
		}

		token, err := utils.GenerateToken(user.Email, user.ID)

		if err != nil {
			fmt.Println(err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not created token with google login."})
			return
		}

		context.JSON(http.StatusCreated, gin.H{"message": "Save user from google info.", "token": token})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not created token with google login."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Google login successfully!", "token": token})
}
