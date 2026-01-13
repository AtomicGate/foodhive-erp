package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService interface {
	UploadFile(bucketName, fileName string, fileContent []byte) (string, error)
	DownloadFile(bucketName, fileName string) ([]byte, error)
	DeleteFile(bucketName, fileName string) error
}

type MinioStorageService struct {
	Client *minio.Client
}

func New(endpoint, accessKey, secretKey string, useSSL bool) (*MinioStorageService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	return &MinioStorageService{Client: client}, nil
}

func (s *MinioStorageService) UploadFile(bucketName, fileName string, fileContent []byte) (string, error) {
	reader := bytes.NewReader(fileContent)
	_, err := s.Client.PutObject(context.Background(), bucketName, fileName, reader, int64(reader.Len()), minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	fmt.Printf("Uploaded to %s: %s/%s\n", s.Client.EndpointURL(), bucketName, fileName)

	return fmt.Sprintf("%v/%v", bucketName, fileName), nil
}

func (s *MinioStorageService) DownloadFile(bucketName, fileName string) ([]byte, error) {
	object, err := s.Client.GetObject(context.Background(), bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer object.Close()

	fileContent, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %v", err)
	}

	return fileContent, nil
}

func (s *MinioStorageService) DeleteFile(bucketName, fileName string) error {
	err := s.Client.RemoveObject(context.Background(), bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
