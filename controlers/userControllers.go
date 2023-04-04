package controlers

import (
	"API-Books/initializer"
	"API-Books/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(ctx *gin.Context) {
	// Get Email

	var body struct {
		Email    string
		Password string
	}

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}
	// Hash The Password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash a Password",
		})
		return
	}
	// Create the user
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
	}
	result := initializer.DB.Create(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	// respose

	ctx.JSON(http.StatusOK, gin.H{})

}

func Login(ctx *gin.Context) {
	// Get Email and Pass off req body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}
	// Look up requested user
	var user models.User
	initializer.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invailid EMAil or Password",
		})
		return
	}

	// Compare sent in pass with saved user Pass Hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invailid Email or PASSword",
		})
		return
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign and get the complete encoded token as stng useing the secret code
	secret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secret))

	//	fmt.Println(time.Now().Add(time.Hour * 24 * 30).Unix())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to create Token",
		})
		return
	}

	// send it Back
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{})
}

func Validate(ctx *gin.Context) {

	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
