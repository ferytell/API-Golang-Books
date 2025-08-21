package controlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath" 
	"log"            


	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// UploadToMinio uploads a file to MinIO and returns its URL
func UploadToMinio(imagePath string) (string, error) {
	   endpoint := os.Getenv("MINIO_ENDPOINT")     // 
	accessKey := os.Getenv("MINIO_ACCESS_KEY") // 
	secretKey := os.Getenv("MINIO_SECRET_KEY") //
	bucketName := os.Getenv("MINIO_BUCKET")    // 
	useSSL := true                            // change to true if using https

	// Initialize minio client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create minio client: %w", err)
	}

	// Ensure bucket exists (create if not exists)
	ctx := context.Background()
	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	 if errBucketExists != nil {
		return "", fmt.Errorf("failed to check bucket: %w", errBucketExists)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Upload the file
	fileName := filepath.Base(imagePath)
	_, err = minioClient.FPutObject(ctx, bucketName, fileName, imagePath, minio.PutObjectOptions{
		ContentType: "image/jpeg", // adjust based on file type
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Build file URL (depends on your MinIO setup)
	url := fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, fileName)

	log.Printf("Successfully uploaded %s\n", fileName)
	return url, nil
}

// func UploadToFlickr(imagePath string) (string, error) {
// 	// Create a new Flickr client with your API key and secret
// 	client := flickr.NewFlickrClient(os.Getenv("FLICKR_API_KEY"), os.Getenv("FLICKR_API_SECRET"))

// 	// Authenticate with Flickr
// 	err := client.Authenticate(flickr.PermissionWrite)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to authenticate with Flickr: %w", err)
// 	}

// 	// Upload the image to Flickr
// 	photoID, err := client.UploadPhoto(context.Background(), imagePath, "")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to upload photo to Flickr: %w", err)
// 	}

// 	// Get the URL of the uploaded photo
// 	photoInfo, err := client.GetPhotoInfo(context.Background(), photoID)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get photo info from Flickr: %w", err)
// 	}
// 	url := photoInfo.URL()

// 	return url, nil
// }
