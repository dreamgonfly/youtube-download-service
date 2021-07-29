package videodownload

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type YoutubeDl struct{}

func (y *YoutubeDl) Preview(id string, dir string) (description string, info string, thumbnail string, err error) {
	args := []string{
		"--skip-download",
		// Output template: https://github.com/ytdl-org/youtube-dl/blob/master/README.md#output-template
		"--output", filepath.Join(dir, "%(title)s_%(format_note)s.%(ext)s"),
		"--write-description",
		"--write-info-json",
		"--write-annotations",
		"--write-sub",
		"--write-thumbnail",
		id,
	}
	_, err = exec.Command("youtube-dl", args...).Output()
	if err != nil {
		return "", "", "", errors.Wrap(err, "could not download")
	}

	stdout, err := exec.Command(
		"youtube-dl",
		"--output", filepath.Join(dir, "%(title)s_%(format_note)s.%(ext)s"),
		"--get-filename",
		id,
	).Output()
	if err != nil {
		return "", "", "", errors.Wrap(err, "could not get filename")
	}
	filename := strings.TrimSpace(string(stdout))
	ext := filepath.Ext(filename)
	basename := filename[:len(filename)-len(ext)] // filename except extention
	description = strings.Join([]string{basename, ".description"}, "")
	info = strings.Join([]string{basename, ".info.json"}, "")
	thumbnail = strings.Join([]string{basename, ".webp"}, "")
	return description, info, thumbnail, nil
}

func (y *YoutubeDl) Download(id string, format_code string, dir string) (video string, err error) {
	args := []string{
		"--format",
		format_code,
		"--output", filepath.Join(dir, "%(title)s_%(format_note)s.%(ext)s"),
		id,
	}
	_, err = exec.Command("youtube-dl", args...).Output()
	if err != nil {
		return "", errors.Wrap(err, "could not download")
	}

	stdout, err := exec.Command(
		"youtube-dl",
		"--format",
		format_code,
		"--output", filepath.Join(dir, "%(title)s_%(format_note)s.%(ext)s"),
		"--get-filename",
		id,
	).Output()
	if err != nil {
		return "", errors.Wrap(err, "could not get filename")
	}
	filename := strings.TrimSpace(string(stdout))
	return filename, nil
}
