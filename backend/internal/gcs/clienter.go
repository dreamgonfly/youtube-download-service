package gcs

import (
	"context"
	"io"
)

type Clienter interface {
	Close()
	Bucket(name string) BucketHandler
}

type BucketHandler interface {
	Object(name string) ObjectHandler
}

type ObjectHandler interface {
	NewReader(ctx context.Context) (io.ReadCloser, error)
	NewWriter(ctx context.Context) io.WriteCloser
}
