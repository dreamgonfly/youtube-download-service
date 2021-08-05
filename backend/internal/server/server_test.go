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
	expected_formats := `{"Thumbnail":"https://storage.googleapis.com/youtube-download-backend-beta/videos/GSVsfCCtRr0/%5B%EA%B8%B0%EC%83%9D%EC%B6%A9%5D%2030%EC%B4%88%20%EC%98%88%EA%B3%A0.jpg?X-Goog-Algorithm=GOOG4-RSA-SHA256\u0026X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210801%2Fauto%2Fstorage%2Fgoog4_request\u0026X-Goog-Date=20210801T084013Z\u0026X-Goog-Expires=899\u0026X-Goog-Signature=31390355ce8353e169279bc4286078c29ef2ed5bcb2bcc13380c67efac26e4155275e0d1a3683cc77e7d9188e33ffaee0e7056c584951e339c327934e7d08a78a0e549acd05b1506ed62c6e048c8500416cc7086f07bd9a03fa1df6b678220b0a1c810d0715f904c889a80fc3208a997ebb6e117f6857cf526daa32bb515c048384f85b27456c6c39ae7647efe2f0be24f0905a52fe0e5b94bbd6871a71c34951f1317cbccee23fa4ed0566415e9b55221fba05558ebdbc9c6e48f11396e1766f7f62b49756b8daa7cc54c6169176f81e6726e91e02dec30d671d7cb4cca85cc6556c5ef312293ece4b22792b53aa05481e3ad0ba225ca3435e935adb25a7143\u0026X-Goog-SignedHeaders=host","Formats":[{"Filesize":1348634,"FormatId":"18","FormatNote":"360p","Ext":"mp4"},{"Filesize":2059470,"FormatId":"22","FormatNote":"720p","Ext":"mp4"}]}`
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
	expected_formats := `{"Thumbnail":"https://storage.googleapis.com/youtube-download-backend-beta/videos/-BIDXOp6_LA/Go%20Modules%20-%20Dependency%20Management%20the%20Right%20Way.webp?X-Goog-Algorithm=GOOG4-RSA-SHA256\u0026X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210801%2Fauto%2Fstorage%2Fgoog4_request\u0026X-Goog-Date=20210801T201751Z\u0026X-Goog-Expires=899\u0026X-Goog-Signature=6baaea933a08dc902f21a4b95f5d88b7fa42664d0ffd83385a24b180c248ab5751c8cc9a9f927517f29235751fef7606825210340e7f995fb4d1609a9c78d153062cb7fa11f67e082c28f50f6262632c3337bc225584b4c15405acfb4b6a03e4d253db41b14d39113bce36140c4afae634a8a9e51dfd08f54700c1512996857dfe6604ecb335228e4baecce9458b160537f97d5f2900448a1edf2da3d2c57da2db7690b8d8c7108762cdf4123ee4cb718352859a0181879bb7cd38ba4b5de679a7fa79ad0ac097819af5910dd2356ee10df14ba3653d9f854f9f4b93778d0d0b6efd547d9070a962e12577055144f0696c4f56c9f2a44c7cbc83fa07e089587b\u0026X-Goog-SignedHeaders=host","Formats":[{"Filesize":75134269,"FormatId":"18","FormatNote":"360p","Ext":"mp4"},{"Filesize":210958377,"FormatId":"22","FormatNote":"720p","Ext":"mp4"}]}`
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
	var jsonStr = []byte(`{"URL":"https://i.ytimg.com/vi/GSVsfCCtRr0/hqdefault.jpg?sqp=-oaymwEcCNACELwBSFXyq4qpAw4IARUAAIhCGAFwAcABBg==&rs=AOn4CLCJj9t8x2PEsTiw4J3l8_nz6kRv0A", "Name": "[기생충] 30초 예고"}`)
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
