package youtubefile

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type YoutubeDl struct{ ExecCommand Commander }

func (y *YoutubeDl) GetName(id string) (name string, err error) {
	stdout, err := y.ExecCommand(
		"youtube-dl",
		"--output", "%(title)s.%(ext)s",
		"--get-filename",
		id,
	).Output()
	if err != nil {
		return "", errors.Wrap(err, "could not get filename")
	}
	filename := strings.TrimSpace(string(stdout))
	ext := filepath.Ext(filename)
	basename := filename[:len(filename)-len(ext)] // filename except extention
	return basename, nil
}

func (y *YoutubeDl) Preview(id, name, dir string) (description string, info string, err error) {
	args := []string{
		"--skip-download",
		// Output template: https://github.com/ytdl-org/youtube-dl/blob/master/README.md#output-template
		"--output", filepath.Join(dir, "%(title)s.%(ext)s"),
		"--write-description",
		"--write-info-json",
		"--write-annotations",
		"--write-sub",
		id,
	}
	_, err = y.ExecCommand("youtube-dl", args...).Output()
	if err != nil {
		return "", "", errors.Wrap(err, "could not download")
	}
	description = filepath.Join(dir, strings.Join([]string{name, ".description"}, ""))
	info = filepath.Join(dir, strings.Join([]string{name, ".info.json"}, ""))
	return description, info, nil
}

func (y *YoutubeDl) Download(id string, format string, dir string) (video string, err error) {
	stdout, err := y.ExecCommand(
		"youtube-dl",
		"--format",
		format,
		"--output", filepath.Join(dir, "%(title)s_%(format_note)s.%(ext)s"),
		"--get-filename",
		id,
	).Output()
	if err != nil {
		return "", errors.Wrap(err, "could not get filename")
	}
	filename := strings.TrimSpace(string(stdout))

	args := []string{
		"--format",
		format,
		"--output", filepath.Join(dir, "%(title)s_%(format_note)s.%(ext)s"),
		id,
	}
	_, err = y.ExecCommand("youtube-dl", args...).Output()
	if err != nil {
		return "", errors.Wrap(err, "could not download")
	}

	return filename, nil
}
