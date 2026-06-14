package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/BounkhongDev/bkgo/config"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Storage is the MinIO adapter implementing contract.Storage.
type Storage struct {
	client        *miniogo.Client
	defaultBucket string
}

// New creates a MinIO client and ensures the default bucket exists.
func New(ctx context.Context, cfg config.MinIO) (*Storage, error) {
	client, err := miniogo.New(cfg.Endpoint, &miniogo.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: connect: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("minio: bucket check: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, miniogo.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio: create bucket: %w", err)
		}
	}

	return &Storage{client: client, defaultBucket: cfg.Bucket}, nil
}

// Client exposes the underlying minio.Client for advanced usage.
func (s *Storage) Client() *miniogo.Client { return s.client }

func (s *Storage) Upload(ctx context.Context, bucket, name string, reader io.Reader, size int64, contentType string) (string, error) {
	if bucket == "" {
		bucket = s.defaultBucket
	}
	_, err := s.client.PutObject(ctx, bucket, name, reader, size, miniogo.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio: upload: %w", err)
	}
	return fmt.Sprintf("%s/%s", bucket, name), nil
}

func (s *Storage) Download(ctx context.Context, bucket, name string) (io.ReadCloser, error) {
	if bucket == "" {
		bucket = s.defaultBucket
	}
	obj, err := s.client.GetObject(ctx, bucket, name, miniogo.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("minio: download: %w", err)
	}
	return obj, nil
}

func (s *Storage) Delete(ctx context.Context, bucket, name string) error {
	if bucket == "" {
		bucket = s.defaultBucket
	}
	return s.client.RemoveObject(ctx, bucket, name, miniogo.RemoveObjectOptions{})
}

func (s *Storage) URL(ctx context.Context, bucket, name string, expiry time.Duration) (string, error) {
	if bucket == "" {
		bucket = s.defaultBucket
	}
	u, err := s.client.PresignedGetObject(ctx, bucket, name, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("minio: presign: %w", err)
	}
	return u.String(), nil
}
