package storage_test

import (
	"testing"
	"youtube-download-backend/storage"

	"github.com/stretchr/testify/assert"
)

func TestNewPathFromURI(t *testing.T) {
	scheme, bucket, key, err := storage.ParseURI("gs://my-bucket/path/to/my/object.json")
	if err != nil {
		t.Fatal("could not make new storage path")
	}
	assert.Equal(t, "gs", scheme, "scheme mismatch")
	assert.Equal(t, "my-bucket", bucket, "bucket mismatch")
	assert.Equal(t, "path/to/my/object.json", key, "key mismatch")
}

func TestNewPathFromURIErr(t *testing.T) {
	scheme, bucket, key, err := storage.ParseURI("")
	assert.NotNil(t, err, "err should be raised")
	assert.Equal(t, "", scheme, "scheme should be empty")
	assert.Equal(t, "", bucket, "bucket should be empty")
	assert.Equal(t, "", key, "key should be empty")
}

func TestNewPathFromKey(t *testing.T) {
	uri := storage.ComposeURI("gs", "my-bucket", "path/to/my/object.json")

	assert.Equal(t, "gs://my-bucket/path/to/my/object.json", uri, "uri mismatch")
}
