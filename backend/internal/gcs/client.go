package gcs

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

type Client struct{ Client *storage.Client }

func (c *Client) Close() {}
func (c *Client) Bucket(name string) BucketHandler {
	return &BucketHandle{c.Client.Bucket(name)}
}

type BucketHandle struct{ BucketHandle *storage.BucketHandle }

func (b *BucketHandle) Object(name string) ObjectHandler {
	return &ObjectHandle{b.BucketHandle.Object(name)}
}

type ObjectHandle struct{ ObjectHandle *storage.ObjectHandle }

func (o *ObjectHandle) NewReader(ctx context.Context) (io.ReadCloser, error) {
	return o.ObjectHandle.NewReader(ctx)
}

func (o *ObjectHandle) NewWriter(ctx context.Context) io.WriteCloser {
	return o.ObjectHandle.NewWriter(ctx)
}
