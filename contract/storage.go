package contract

import (
	"context"
	"io"
	"time"
)

// Storage is the port for any object/file storage adapter.
// Implement this interface to swap between MinIO, AWS S3, GCS, local disk, etc.
type Storage interface {
	Upload(ctx context.Context, bucket, name string, reader io.Reader, size int64, contentType string) (string, error)
	Download(ctx context.Context, bucket, name string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, name string) error
	// URL returns a pre-signed URL valid for the given expiry duration.
	URL(ctx context.Context, bucket, name string, expiry time.Duration) (string, error)
}
