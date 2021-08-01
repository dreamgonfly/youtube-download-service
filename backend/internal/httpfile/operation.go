package httpfile

import (
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type Clienter interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}

// https://golangcode.com/download-a-file-from-a-url/
func DownloadFile(client Clienter, url string, path string) error {
	resp, err := client.Get(url)
	if err != nil {
		return errors.Wrapf(err, "failed to get url (%s)", url)
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to create file (%s)", path)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
