package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/yourusername/oden/internal/config"
)

// Client represents a storage client
type Client struct {
	client     *minio.Client
	bucketName string
}

// NewClient creates a new storage client
func NewClient(cfg *config.Config) (*Client, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.Storage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Storage.AccessKey, cfg.Storage.SecretKey, ""),
		Secure: cfg.Storage.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating MinIO client: %w", err)
	}

	// Check if bucket exists
	exists, err := minioClient.BucketExists(context.Background(), cfg.Storage.Bucket)
	if err != nil {
		return nil, fmt.Errorf("error checking bucket existence: %w", err)
	}

	// If the bucket doesn't exist and we're not in production, create it
	if !exists {
		err = minioClient.MakeBucket(context.Background(), cfg.Storage.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("error creating bucket: %w", err)
		}

		// Set bucket policy to allow public read access
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": "*",
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + cfg.Storage.Bucket + `/*"]
				}
			]
		}`
		err = minioClient.SetBucketPolicy(context.Background(), cfg.Storage.Bucket, policy)
		if err != nil {
			return nil, fmt.Errorf("error setting bucket policy: %w", err)
		}
	}

	return &Client{
		client:     minioClient,
		bucketName: cfg.Storage.Bucket,
	}, nil
}

// UploadFile uploads a file to the storage
func (c *Client) UploadFile(objectName string, filePath string, contentType string) error {
	_, err := c.client.FPutObject(context.Background(), c.bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}
	return nil
}

// GetFileURL gets the URL of a file in the storage
func (c *Client) GetFileURL(objectName string) string {
	return fmt.Sprintf("http://%s/%s/%s", c.client.EndpointURL().Host, c.bucketName, objectName)
}

// DeleteFile deletes a file from the storage
func (c *Client) DeleteFile(objectName string) error {
	err := c.client.RemoveObject(context.Background(), c.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}
	return nil
} 