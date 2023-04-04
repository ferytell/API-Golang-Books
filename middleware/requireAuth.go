package middleware

import (
	"API-Books/initializer"
	"API-Books/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(ctx *gin.Context) {

	// Get Cookies from req
	//	fmt.Println("in Middleware")
	tokenString, err := ctx.Cookie("Authorization")

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpexted signing method: %v", token.Header["alg"])
		}

		// HmacSampleSecret is a [byte] containing your secret, e.g []byte("my_secret")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check expire

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		// find user with token

		var user models.User
		initializer.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		// attach it to the request

		ctx.Set("User", user)

		ctx.Next()

	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

}
