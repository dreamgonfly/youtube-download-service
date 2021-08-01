package gcs

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"time"
	"youtube-download-backend/internal/config"

	"github.com/pkg/errors"
)

// https://cloud.google.com/storage/docs/downloading-objects#storage-download-object-go
func DownloadFile(ctx context.Context, client Clienter, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Hour)
	defer cancel()

	rc, err := client.Bucket(config.Conf.Bucket).Object(key).NewReader(ctx)
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
func UploadFile(ctx context.Context, client Clienter, path, key string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "could not open path")
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Hour)
	defer cancel()

	wc := client.Bucket(config.Conf.Bucket).Object(key).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return errors.Wrap(err, "io.Copy")
	}
	if err := wc.Close(); err != nil {
		return errors.Wrap(err, "Writer.Close")
	}
	return nil
}
