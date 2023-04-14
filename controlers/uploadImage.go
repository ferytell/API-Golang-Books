package controlers

import (
	"context"
	"net/http"
	"os"

	"github.com/codedius/imagekit-go"
	"github.com/gin-gonic/gin"
)

func UploadImage(ctx *gin.Context) {
	// Replace with your own API keys
	publicKey := os.Getenv("IMAGEKITPUBLICKEY")
	privateKey := os.Getenv("IMAGEKITPRIVATEKEY")

	opts := imagekit.Options{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	// Create a new ImageKit client
	client, err := imagekit.NewClient(&opts)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "fail to create client imege",
		})
	}

	// Open the image file to upload
	file, err := os.Open("image.jpg")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to open image file",
		})
		return
	}
	defer file.Close()

	// Upload the image to ImageKit
	uploadParams := imagekit.UploadRequest{
		File:              file,
		FileName:          "image.jpg",
		UseUniqueFileName: false,
		Tags:              []string{"go", "image"},
		Folder:            "/",
		IsPrivateFile:     false,
	}

	c := context.Background()
	uploadResult, err := client.Upload.ServerUpload(c, &uploadParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to upload image",
		})
		return
	}
	// Return the URL of the uploaded image as a JSON response
	ctx.JSON(http.StatusOK, gin.H{
		"url": uploadResult.URL,
	})
}

// import (
// 	"fmt"
// 	"os"

// 	"github.com/koffeinsource/go-imgur/imgur"
// )

// func UploadImage() {
// 	// Load the client ID and secret from environment variables
// 	clientID := os.Getenv("IMGUR_CLIENT_ID")
// 	clientSecret := os.Getenv("IMGUR_CLIENT_SECRET")

// 	// Create a new Imgur client
// 	client, err := imgur.NewClient(&imgur.Auth{
// 		ClientID:     clientID,
// 		ClientSecret: clientSecret,
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// Open the image file to upload
// 	file, err := os.Open("image.jpg")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer file.Close()

// 	// Upload the image to Imgur
// 	response, err := client.UploadFromReader(file)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// Print the URL of the uploaded image
// 	fmt.Println(response.Link)
// }
