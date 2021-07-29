package videodownload

import (
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
)

type DownloaderMock struct{}

func (d *DownloaderMock) Preview(id string, dir string) (description string, info string, thumbnail string, err error) {
	if id == "x5TLTSGrn_M" {
		descname := "‘교도소 다녀오면 5억 줄게’…치밀한 범행 계획 _ KBS 2021.05.14.-x5TLTSGrn_M.description"
		infoname := "‘교도소 다녀오면 5억 줄게’…치밀한 범행 계획 _ KBS 2021.05.14.-x5TLTSGrn_M.info.json"
		thumbnailname := "‘교도소 다녀오면 5억 줄게’…치밀한 범행 계획 _ KBS 2021.05.14.-x5TLTSGrn_M.webp"
		err = copyFileContents(path.Join("../testdata", descname), path.Join(dir, descname))
		if err != nil {
			return "", "", "", errors.Wrap(err, "could not copy file contents")
		}
		err = copyFileContents(path.Join("../testdata", infoname), path.Join(dir, infoname))
		if err != nil {
			return "", "", "", errors.Wrap(err, "could not copy file contents")
		}
		err = copyFileContents(path.Join("../testdata", thumbnailname), path.Join(dir, thumbnailname))
		if err != nil {
			return "", "", "", errors.Wrap(err, "could not copy file contents")
		}
		return path.Join(dir, descname), path.Join(dir, infoname), path.Join(dir, thumbnailname), nil
	} else {
		return "", "", "", errors.New("unable to mock")
	}
}

func (d *DownloaderMock) Download(id string, format_code string, dir string) (video string, err error) {
	if id == "x5TLTSGrn_M" && format_code == "22" {
		videoname := "‘교도소 다녀오면 5억 줄게’…치밀한 범행 계획 _ KBS 2021.05.14.-x5TLTSGrn_M.mp4"
		err := copyFileContents(path.Join("../testdata", videoname), path.Join(dir, videoname))
		if err != nil {
			return "", errors.Wrap(err, "could not copy file contents")
		}
		return path.Join(dir, videoname), nil
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
