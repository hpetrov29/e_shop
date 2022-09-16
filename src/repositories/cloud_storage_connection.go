package repositories

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/fnmzgdt/e_shop/src/utils"
)

func SetupGoogleStorageConnection() (*storage.Client, error) {
	fmt.Println("Connecting to Cloud Storage")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	storage, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error while initializing Cloud Storage Client: %w", err)
	}
	return storage, nil
}

func (bucket GCBucket) CloseConnection() error {
	if err := bucket.Client.Close(); err != nil {
		return fmt.Errorf("Error while closing Google Storage Client: %w", err)
	}
	return nil
}

type GCBucket struct {
	Client     *storage.Client
	BucketName string
}

func NewGCBucket(gcClient *storage.Client, bucketName string) *GCBucket {
	return &GCBucket{Client: gcClient, BucketName: bucketName}
}

func (b *GCBucket) UploadImage(ctx context.Context, objName string, imageFile multipart.File) (string, error) {
	bucket := b.Client.Bucket(b.BucketName)
	object := bucket.Object(objName)
	writer := object.NewWriter(ctx)
	if _, err := io.Copy(writer, imageFile); err != nil {
		log.Printf("Unable to write file to Google Cloud Storage: %v", err)
		return "", nil
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}
	bucketUrl := utils.GetEnv("PUBLIC_URL", "")
	imageUrl := fmt.Sprintf(bucketUrl, b.BucketName, objName)
	return imageUrl, nil
}

func (b *GCBucket) DeleteImage(ctx context.Context, objName string) error {
	bucket := b.Client.Bucket(b.BucketName)
	object := bucket.Object(objName)
	if err := object.Delete(ctx); err != nil {
		return err
	}
	return nil
}
