package gcs

import (
	"os"
	"time"
	. "youtube-download-backend/internal/config"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
)

type SignFunc func(bucket, name string, opts *storage.SignedURLOptions) (string, error)

// https://cloud.google.com/storage/docs/access-control/signing-urls-with-helpers#storage-signed-url-object-go
func GenerateV4GetObjectSignedURL(sign SignFunc, key string, duration time.Duration) (string, error) {
	serviceAccount := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	jsonKey, err := os.ReadFile(serviceAccount)
	if err != nil {
		return "", errors.Wrap(err, "ioutil.ReadFile")
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", errors.Wrap(err, "google.JWTConfigFromJSON")
	}
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(duration),
	}
	u, err := sign(Config.Bucket, key, opts)
	if err != nil {
		return "", errors.Wrap(err, "storage.SignedURL")
	}
	return u, nil
}
