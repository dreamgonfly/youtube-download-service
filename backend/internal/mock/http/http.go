package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"youtube-download-backend/internal/config"

	"github.com/pkg/errors"
)

type Client struct {
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not create mock request %s", url))
	}
	return c.Do(req)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Path == "/vi/GSVsfCCtRr0/hqdefault.jpg" {
		f, err := os.Open(filepath.Join(config.RootDir, "testdata", "[기생충] 30초 예고.jpg"))
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: 200,
			Body:       f,
		}, nil
	} else if req.URL.Path == "/vi_webp/-BIDXOp6_LA/maxresdefault.webp" {
		f, err := os.Open(filepath.Join(config.RootDir, "testdata", "Go Modules - Dependency Management the Right Way.webp"))
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: 200,
			Body:       f,
		}, nil
	} else {
		return nil, errors.New(fmt.Sprintf("could not mock %s", req.URL.Path))
	}
}
