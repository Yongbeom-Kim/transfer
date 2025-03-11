package storage

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var client *storage.Client = nil

func Client() (*storage.Client, error) {
	if client != nil {
		return client, nil
	}
	var err error
	if os.Getenv("BACKEND_GOOGLE_APPLICATION_CREDENTIALS") == "" {
		return nil, errors.New("BACKEND_GOOGLE_APPLICATION_CREDENTIALS is not set")
	}
	client, err = storage.NewClient(context.Background(), option.WithCredentialsFile(os.Getenv("BACKEND_GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Bucket() (*storage.BucketHandle, error) {
	bucketName := os.Getenv("GCLOUD_BUCKET_NAME")
	if bucketName == "" {
		return nil, errors.New("GCLOUD_BUCKET_NAME is not set")
	}
	client, err := Client()
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketName), nil
}
