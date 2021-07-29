package videodownload

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type DownloaderMock struct{}

func (d *DownloaderMock) GetName(id string) (name string, err error) {
	return "[기생충] 30초 예고", nil
}

func (d *DownloaderMock) Preview(id, name, dir string) (description string, info string, err error) {
	if id == "GSVsfCCtRr0" {
		descname := "[기생충] 30초 예고.description"
		infoname := "[기생충] 30초 예고.info.json"
		err = copyFileContents(filepath.Join("../testdata", descname), filepath.Join(dir, descname))
		if err != nil {
			return "", "", errors.Wrap(err, "could not copy file contents")
		}
		err = copyFileContents(filepath.Join("../testdata", infoname), filepath.Join(dir, infoname))
		if err != nil {
			return "", "", errors.Wrap(err, "could not copy file contents")
		}
		return filepath.Join(dir, descname), filepath.Join(dir, infoname), nil
	} else {
		return "", "", errors.New("unable to mock")
	}
}

func (d *DownloaderMock) Download(id string, format_code string, dir string) (video string, err error) {
	if id == "GSVsfCCtRr0" && format_code == "18" {
		videoname := "[기생충] 30초 예고_360p.mp4"
		err := copyFileContents(filepath.Join("../testdata", videoname), filepath.Join(dir, videoname))
		if err != nil {
			return "", errors.Wrap(err, "could not copy file contents")
		}
		return filepath.Join(dir, videoname), nil
	} else {
		return "", errors.New("unable to mock")
	}
}

// https://stackoverflow.com/a/21067803/7866795
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	if err != nil {
		return err
	}
	return
}
