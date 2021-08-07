package youtubefile

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
		"--", // The '--' tells the shell that what follows after this is not an option to the command.
		id,
	).Output()
	if err != nil {
		return "", errors.Wrap(err, "could not get filename")
	}
	filename := strings.TrimSpace(string(stdout))
	return filename, nil
}

func (y *YoutubeDl) GetFilenameWithFormat(id, format string) (string, error) {
	cmd := y.ExecCommand(
		"youtube-dl",
		"--format",
		format,
		"--output", "%(title)s_%(format_note)s.%(ext)s",
		"--get-filename",
		"--",
		id,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("stdout error command (%s)", cmd.String()))
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("stderr error command (%s)", cmd.String()))
	}

	err = cmd.Start()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("error starting (%s)", cmd.String()))
		return "", err
	}
	out, err := io.ReadAll(stdout)
	if err != nil {
		return "", errors.Wrap(err, "could not read stdout")
	}
	errout, err := io.ReadAll(stderr)
	if err != nil {
		err = errors.Wrap(err, "could not read stderr")
	}
	err = cmd.Wait()
	if err != nil {
		err = errors.Wrap(err, strings.TrimSpace(string(errout)))
		return "", errors.Wrap(err, fmt.Sprintf("error waiting command (%s)", cmd.String()))
	}

	filename := strings.TrimSpace(string(out))
	return filename, nil
}

func (y *YoutubeDl) Preview(id, dir string) (description string, info string, err error) {
	args := []string{
		"--skip-download",
		// Output template: https://github.com/ytdl-org/youtube-dl/blob/master/README.md#output-template
		"--output", filepath.Join(dir, "%(title)s.%(ext)s"),
		"--write-description",
		"--write-info-json",
		"--write-annotations",
		"--write-sub",
		"--",
		id,
	}
	out, err := y.ExecCommand("youtube-dl", args...).CombinedOutput()
	if err != nil {
		return "", "", errors.Wrap(err, fmt.Sprintf("error excuting (youtube-dl %s) outputing (%s)", strings.Join(args, " "), string(out)))
	}
	name, err := GetNameFromDescription(dir)
	if err != nil {
		return "", "", errors.Wrap(err, "could not get name from description")
	}
	description = filepath.Join(dir, strings.Join([]string{name, ".description"}, ""))
	info = filepath.Join(dir, strings.Join([]string{name, ".info.json"}, ""))
	return description, info, nil
}

func (y *YoutubeDl) Download(id string, format string, dir string) (video string, err error) {
	args := []string{
		"--format",
		format,
		"--output", filepath.Join(dir, "%(title)s_%(format_note)s.%(ext)s"),
		"--",
		id,
	}
	out, err := y.ExecCommand("youtube-dl", args...).CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("error excuting (youtube-dl %s) outputing (%s)", strings.Join(args, " "), string(out)))
	}
	name, err := GetVideoPathFromDir(dir)
	if err != nil {
		return "", errors.Wrap(err, "could not get name from dir")
	}

	return filepath.Join(dir, name), nil
}

func (y *YoutubeDl) StreamDownloadCommand(id string, format string, w http.ResponseWriter) (cmd Outputer) {

	args := []string{
		"--format",
		format,
		"--force-ipv4", // Youtube ban entire IPv6 blocks. https://github.com/ytdl-org/youtube-dl/issues/21729
		"--output", "-",
		"--",
		id,
	}
	cmd = y.ExecCommand("youtube-dl", args...)
	return cmd
}

func Stem(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)] // filename except extention
}

// GetNameFromDescription assumes there is only one description file in dir
func GetNameFromDescription(dir string) (name string, err error) {
	ext := ".description"
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", errors.Wrap(err, "could not read dir")
	}
	for _, f := range files {
		n := f.Name()
		if strings.HasSuffix(n, ext) {
			name = Stem(n)
			return name, nil
		}
	}
	return "", errors.New(fmt.Sprintf("No file with %s extension", ext))
}

// GetVideoPathFromDir assumes there is only one file in dir
func GetVideoPathFromDir(dir string) (name string, err error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", errors.Wrap(err, "could not read dir")
	}
	if len(files) != 1 {
		return "", errors.Wrap(err, fmt.Sprintf("expected 1 file. got %d files", len(files)))
	}
	f := files[0]
	return f.Name(), nil
}
