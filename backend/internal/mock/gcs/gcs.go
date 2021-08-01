package gcs

import (
	"bytes"
	"context"
	"io"
	"youtube-download-backend/internal/gcs"
)

type Client struct{}

func (c *Client) Close() {}
func (c *Client) Bucket(name string) gcs.BucketHandler {
	return &BucketHandle{}
}

type BucketHandle struct{}

func (b *BucketHandle) Object(name string) gcs.ObjectHandler {
	return &ObjectHandle{}
}

type ObjectHandle struct{}

func (o *ObjectHandle) NewReader(ctx context.Context) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte{})), nil
}

func (o *ObjectHandle) NewWriter(ctx context.Context) io.WriteCloser {
	return &DummyWriteCloser{}
}

type DummyWriteCloser struct{}

func (d *DummyWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (d *DummyWriteCloser) Close() error {
	return nil
}
