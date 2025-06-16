package minio

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"os"
)

var Client *minio.Client
var BucketName string

func InitMinio() {
	endpoint := os.Getenv("MINIO_ENDPOINT")    // ex: "minio:9000"
	accessKey := os.Getenv("MINIO_ACCESS_KEY") // ex: "minioadmin"
	secretKey := os.Getenv("MINIO_SECRET_KEY") // ex: "minioadmin"
	BucketName = os.Getenv("MINIO_BUCKET")     // ex: "attachments"
	useSSL := false

	// Initialize client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("❌ MinIO client init failed: %v", err)
	}

	Client = minioClient
	log.Printf("✅ MinIO client initialized — Bucket: %s", BucketName)

	// Check if bucket exists
	exists, err := Client.BucketExists(context.Background(), BucketName)
	if err != nil {
		log.Fatalf("❌ Error checking bucket: %v", err)
	}
	if !exists {
		log.Fatalf("🚫 Bucket '%s' does not exist. Create it manually first.", BucketName)
	}
}
func UploadFile(originalName string, file io.Reader) (string, string, error) {
	hashedName := uuid.New().String()

	uploadInfo, err := Client.PutObject(
		context.Background(),
		BucketName,
		hashedName,
		file,
		-1,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		return "", "", err
	}
	log.Printf("📦 Uploaded %s (%d bytes)", uploadInfo.Key, uploadInfo.Size)

	url := fmt.Sprintf("http://localhost:9000/%s/%s", BucketName, hashedName)
	return url, hashedName, nil
}

func DownloadFile(hashedName string) (io.ReadCloser, error) {
	object, err := Client.GetObject(context.Background(), BucketName, hashedName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
