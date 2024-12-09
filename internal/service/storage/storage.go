package storage_service

import (
	"account/internal/common"
	"account/internal/infrastructure/configuration"
	"bytes"
	"context"
	"mime"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type service struct {
	s3         *minio.Client
	bucketName string
}

func New(cfg *configuration.S3Config) (*service, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &service{
		s3:         client,
		bucketName: cfg.BucketName,
	}, nil
}

func (s *service) Put(c context.Context, data []byte) (string, error) {
	filename, err := common.GenerateRandomHash()
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(data)

	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", err
	}

	filename = filename + exts[len(exts)-1]

	if _, err := s.s3.PutObject(
		c, s.bucketName, filename, bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{
			ContentType: mimeType,
		},
	); err != nil {
		return "", err
	}

	return filename, nil
}

func (s *service) Remove(c context.Context, filename string) error {
	return s.s3.RemoveObject(c, s.bucketName, filename, minio.RemoveObjectOptions{})
}
