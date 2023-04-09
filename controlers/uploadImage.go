package controlers

import (
	"context"
	"fmt"
	"os"

	"github.com/azer/go-flickr"
)

func UploadToFlickr(imagePath string) (string, error) {
	// Create a new Flickr client with your API key and secret
	client := flickr.NewFlickrClient(os.Getenv("FLICKR_API_KEY"), os.Getenv("FLICKR_API_SECRET"))

	// Authenticate with Flickr
	err := client.Authenticate(flickr.PermissionWrite)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate with Flickr: %w", err)
	}

	// Upload the image to Flickr
	photoID, err := client.UploadPhoto(context.Background(), imagePath, "")
	if err != nil {
		return "", fmt.Errorf("failed to upload photo to Flickr: %w", err)
	}

	// Get the URL of the uploaded photo
	photoInfo, err := client.GetPhotoInfo(context.Background(), photoID)
	if err != nil {
		return "", fmt.Errorf("failed to get photo info from Flickr: %w", err)
	}
	url := photoInfo.URL()

	return url, nil
}
