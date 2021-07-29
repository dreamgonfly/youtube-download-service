package storage

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

type GCS struct{}

// https://cloud.google.com/storage/docs/downloading-objects#storage-download-object-go
func (g *GCS) DownloadFile(uri string) ([]byte, error) {
	scheme, bucket, key, err := ParseURI(uri)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing uri")
	}
	if scheme != "gs" {
		return nil, errors.New("scheme should be gs")
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "storage.NewClient")
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Hour*2)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(key).NewReader(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Object(%q).NewReader", key)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}
	return data, nil
}

// https://cloud.google.com/storage/docs/uploading-objects#storage-upload-object-go
func (g *GCS) UploadFile(path string, uri string) error {
	scheme, bucket, key, err := ParseURI(uri)
	if err != nil {
		return errors.Wrap(err, "error parsing uri")
	}
	if scheme != "gs" {
		return errors.New("scheme should be gs")
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errors.Wrap(err, "storage.NewClient")
	}
	defer client.Close()

	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "could not open path")
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Hour*2)
	defer cancel()

	wc := client.Bucket(bucket).Object(key).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return errors.Wrap(err, "io.Copy")
	}
	if err := wc.Close(); err != nil {
		return errors.Wrap(err, "Writer.Close")
	}
	return nil
}
