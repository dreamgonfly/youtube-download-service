package server_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"youtube-download-backend/server"
	"youtube-download-backend/storage"
	"youtube-download-backend/videodownload"

	"github.com/stretchr/testify/assert"
)

func TestHandleHello(t *testing.T) {
	st := &storage.StorerMock{}
	dw := &videodownload.DownloaderMock{}
	srv := httptest.NewServer(server.NewServer(st, dw))
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
	st := &storage.StorerMock{}
	dw := &videodownload.DownloaderMock{}
	srv := httptest.NewServer(server.NewServer(st, dw))
	url := fmt.Sprintf("%s/preview/x5TLTSGrn_M", srv.URL)
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
	expected_formats := `[{"Filesize":8082004,"FormatId":"18","FormatNote":"360p","Ext":"mp4"},{"Filesize":9404913,"FormatId":"22","FormatNote":"720p","Ext":"mp4"}]`
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
	expected, err := os.ReadFile("../testdata/‘교도소 다녀오면 5억 줄게’…치밀한 범행 계획 _ KBS 2021.05.14.-x5TLTSGrn_M.webp")
	if err != nil {
		t.Fatalf("could not read expected thumbnail: %v", err)
	}
	assert.Equal(t, expected, actual)
}

func TestHandleDownload(t *testing.T) {
	st := &storage.StorerMock{}
	dw := &videodownload.DownloaderMock{}
	srv := httptest.NewServer(server.NewServer(st, dw))
	url := fmt.Sprintf("%s/download/x5TLTSGrn_M?format_code=22", srv.URL)
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
	expected, err := os.ReadFile("../testdata/‘교도소 다녀오면 5억 줄게’…치밀한 범행 계획 _ KBS 2021.05.14.-x5TLTSGrn_M.mp4")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	assert.Equal(t, expected, actual)
}
