package server_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	execmock "youtube-download-backend/internal/mock/exec"
	gcsmock "youtube-download-backend/internal/mock/gcs"
	httpmock "youtube-download-backend/internal/mock/http"
	"youtube-download-backend/internal/server"

	"github.com/stretchr/testify/assert"
)

func TestHandleHello(t *testing.T) {
	ctx := context.Background()
	c := execmock.Command
	g := &gcsmock.Client{}
	h := &httpmock.Client{}
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h))
	res, err := http.Get(fmt.Sprintf("%s/hello", srv.URL))
	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("could not read response: %v", err)
	}

	assert.Equal(t, "Hello!", string(bytes.TrimSpace(b)))
}

func TestHandlePreview(t *testing.T) {
	ctx := context.Background()
	c := execmock.Command
	g := &gcsmock.Client{}
	h := &httpmock.Client{}
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h))
	url := fmt.Sprintf("%s/preview/GSVsfCCtRr0", srv.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}
	req.Header.Set("Accept", "multipart/form-data; charset=utf-8")

	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("could not process request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	_, params, _ := mime.ParseMediaType(res.Header.Get("Content-Type"))
	mr := multipart.NewReader(res.Body, params["boundary"])

	// Formats Part
	part, err := mr.NextPart()
	if err != nil {
		t.Fatalf("could not get NextPart: %v", err)
	}
	value, err := ioutil.ReadAll(part)
	if err != nil {
		t.Fatalf("could not read value: %v", err)
	}
	actual_formats := string(bytes.TrimSpace(value))
	expected_formats := `[{"Filesize":1348634,"FormatId":"18","FormatNote":"360p","Ext":"mp4"},{"Filesize":2059470,"FormatId":"22","FormatNote":"720p","Ext":"mp4"}]`
	assert.Equal(t, expected_formats, actual_formats, "formats mismatch")

	// Thumbnail Part
	part, err = mr.NextPart()
	if err != nil {
		t.Fatalf("could not get NextPart: %v", err)
	}

	actual, err := ioutil.ReadAll(part)
	if err != nil {
		t.Fatalf("could not read actual thumbnail: %v", err)
	}
	expected, err := os.ReadFile("../../testdata/[기생충] 30초 예고.jpg")
	if err != nil {
		t.Fatalf("could not read expected thumbnail: %v", err)
	}
	assert.Equal(t, expected, actual)
}

func TestHandleDownload(t *testing.T) {
	ctx := context.Background()
	c := execmock.Command
	g := &gcsmock.Client{}
	h := &httpmock.Client{}
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h))
	url := fmt.Sprintf("%s/download/GSVsfCCtRr0?format=18", srv.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}
	req.Header.Set("Accept", "multipart/form-data; charset=utf-8")

	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("could not process request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	actual, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
	expected, err := os.ReadFile("../../testdata/[기생충] 30초 예고_360p.mp4")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	assert.Equal(t, expected, actual)
}
