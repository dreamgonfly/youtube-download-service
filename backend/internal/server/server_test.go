package server_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
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
	sf := gcsmock.SignedURL
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h, sf))
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
	sf := gcsmock.SignedURL
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h, sf))
	url := fmt.Sprintf("%s/preview/GSVsfCCtRr0", srv.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("could not process request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	value, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read value: %v", err)
	}
	actual_formats := string(bytes.TrimSpace(value))
	expected_formats := `{"Title":"[기생충] 30초 예고","DurationSecond":30,"Thumbnail":"https://i.ytimg.com/vi/GSVsfCCtRr0/hqdefault.jpg?sqp=-oaymwEcCNACELwBSFXyq4qpAw4IARUAAIhCGAFwAcABBg==\u0026rs=AOn4CLCJj9t8x2PEsTiw4J3l8_nz6kRv0A","Name":"[기생충] 30초 예고","Formats":[{"Filesize":1348634,"FormatId":"18","FormatNote":"360p","Ext":"mp4"},{"Filesize":2059470,"FormatId":"22","FormatNote":"720p","Ext":"mp4"}]}`
	assert.Equal(t, expected_formats, actual_formats, "formats mismatch")
}
func TestHandlePreviewWithDashVideoId(t *testing.T) {
	ctx := context.Background()
	c := execmock.Command
	g := &gcsmock.Client{}
	h := &httpmock.Client{}
	sf := gcsmock.SignedURL
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h, sf))
	url := fmt.Sprintf("%s/preview/-BIDXOp6_LA", srv.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("could not process request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	value, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read value: %v", err)
	}
	actual_formats := string(bytes.TrimSpace(value))
	expected_formats := `{"Title":"Go Modules: Dependency Management the Right Way","DurationSecond":3073,"Thumbnail":"https://i.ytimg.com/vi_webp/-BIDXOp6_LA/maxresdefault.webp","Name":"Go Modules - Dependency Management the Right Way","Formats":[{"Filesize":75134269,"FormatId":"18","FormatNote":"360p","Ext":"mp4"},{"Filesize":210958377,"FormatId":"22","FormatNote":"720p","Ext":"mp4"}]}`
	assert.Equal(t, expected_formats, actual_formats, "formats mismatch")
}

func TestUpdateThumbnail(t *testing.T) {
	ctx := context.Background()
	c := execmock.Command
	g := &gcsmock.Client{}
	h := &httpmock.Client{}
	sf := gcsmock.SignedURL
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h, sf))
	url := fmt.Sprintf("%s/update-thumbnail/GSVsfCCtRr0", srv.URL)
	var jsonStr = []byte(`{"URL":"https://i.ytimg.com/vi/GSVsfCCtRr0/hqdefault.jpg?sqp=-oaymwEcCNACELwBSFXyq4qpAw4IARUAAIhCGAFwAcABBg==\u0026rs=AOn4CLCJj9t8x2PEsTiw4J3l8_nz6kRv0A", "Name": "[기생충] 30초 예고"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("could not create POST request: %v", err)
	}

	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("could not process request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	value, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read value: %v", err)
	}
	actual_formats := string(bytes.TrimSpace(value))
	expected_formats := `{"Thumbnail":"https://storage.googleapis.com/youtube-download-backend-beta/videos/GSVsfCCtRr0/%5B%EA%B8%B0%EC%83%9D%EC%B6%A9%5D%2030%EC%B4%88%20%EC%98%88%EA%B3%A0.jpg?X-Goog-Algorithm=GOOG4-RSA-SHA256\u0026X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210801%2Fauto%2Fstorage%2Fgoog4_request\u0026X-Goog-Date=20210801T084013Z\u0026X-Goog-Expires=899\u0026X-Goog-Signature=31390355ce8353e169279bc4286078c29ef2ed5bcb2bcc13380c67efac26e4155275e0d1a3683cc77e7d9188e33ffaee0e7056c584951e339c327934e7d08a78a0e549acd05b1506ed62c6e048c8500416cc7086f07bd9a03fa1df6b678220b0a1c810d0715f904c889a80fc3208a997ebb6e117f6857cf526daa32bb515c048384f85b27456c6c39ae7647efe2f0be24f0905a52fe0e5b94bbd6871a71c34951f1317cbccee23fa4ed0566415e9b55221fba05558ebdbc9c6e48f11396e1766f7f62b49756b8daa7cc54c6169176f81e6726e91e02dec30d671d7cb4cca85cc6556c5ef312293ece4b22792b53aa05481e3ad0ba225ca3435e935adb25a7143\u0026X-Goog-SignedHeaders=host"}`

	assert.Equal(t, expected_formats, actual_formats, "formats mismatch")
}

func TestHandleDownload(t *testing.T) {
	ctx := context.Background()
	c := execmock.Command
	g := &gcsmock.Client{}
	h := &httpmock.Client{}
	sf := gcsmock.SignedURL
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h, sf))
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

func TestHandleDownloadWithFilename(t *testing.T) {
	ctx := context.Background()
	c := execmock.Command
	g := &gcsmock.Client{}
	h := &httpmock.Client{}
	sf := gcsmock.SignedURL
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h, sf))
	url := fmt.Sprintf("%s/download/GSVsfCCtRr0?format=18&filename=test.mp4", srv.URL)
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

func TestHandlePlay(t *testing.T) {
	ctx := context.Background()
	c := execmock.Command
	g := &gcsmock.Client{}
	h := &httpmock.Client{}
	sf := gcsmock.SignedURL
	srv := httptest.NewServer(server.NewServer(ctx, c, g, h, sf))
	url := fmt.Sprintf("%s/play/GSVsfCCtRr0?format=18", srv.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("could not process request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	value, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read value: %v", err)
	}
	actual_formats := string(bytes.TrimSpace(value))
	expected_formats := `{"Filename":"[기생충] 30초 예고_360p.mp4","URL":"https://storage.googleapis.com/youtube-download-backend-beta/videos/GSVsfCCtRr0/test.mp4?X-Goog-Algorithm=GOOG4-RSA-SHA256\u0026X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210801%2Fauto%2Fstorage%2Fgoog4_request\u0026X-Goog-Date=20210801T084013Z\u0026X-Goog-Expires=899\u0026X-Goog-Signature=31390355ce8353e169279bc4286078c29ef2ed5bcb2bcc13380c67efac26e4155275e0d1a3683cc77e7d9188e33ffaee0e7056c584951e339c327934e7d08a78a0e549acd05b1506ed62c6e048c8500416cc7086f07bd9a03fa1df6b678220b0a1c810d0715f904c889a80fc3208a997ebb6e117f6857cf526daa32bb515c048384f85b27456c6c39ae7647efe2f0be24f0905a52fe0e5b94bbd6871a71c34951f1317cbccee23fa4ed0566415e9b55221fba05558ebdbc9c6e48f11396e1766f7f62b49756b8daa7cc54c6169176f81e6726e91e02dec30d671d7cb4cca85cc6556c5ef312293ece4b22792b53aa05481e3ad0ba225ca3435e935adb25a7143\u0026X-Goog-SignedHeaders=host"}`

	assert.Equal(t, expected_formats, actual_formats, "mismatch")
}
