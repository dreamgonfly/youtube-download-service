package exec

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"youtube-download-backend/internal/config"
	"youtube-download-backend/internal/youtubefile"
)

func Command(name string, arg ...string) youtubefile.Outputer {
	if name == "youtube-dl" &&
		arg[0] == "--output" &&
		arg[1] == "%(title)s.%(ext)s" &&
		arg[2] == "--get-filename" &&
		arg[3] == "--" &&
		arg[4] == "GSVsfCCtRr0" {
		return &Cmd{Content: []byte("[기생충] 30초 예고.mp4")}
	} else if name == "youtube-dl" &&
		arg[0] == "--output" &&
		arg[1] == "%(title)s.%(ext)s" &&
		arg[2] == "--get-filename" &&
		arg[3] == "--" &&
		arg[4] == "-BIDXOp6_LA" {
		return &Cmd{Content: []byte("Go Modules - Dependency Management the Right Way.mp4")}
	} else if name == "youtube-dl" &&
		arg[0] == "--skip-download" &&
		arg[1] == "--output" &&
		arg[2][:4] == "/var" && // /var/folders/53/tpn8zp511y1gdz9k_srhvbh80000gn/T/317835077/%(title)s.%(ext)s
		arg[3] == "--write-description" &&
		arg[4] == "--write-info-json" &&
		arg[5] == "--write-annotations" &&
		arg[6] == "--write-sub" &&
		arg[7] == "--" &&
		arg[8] == "GSVsfCCtRr0" {
		dir := filepath.Dir(arg[2])
		descname := "[기생충] 30초 예고.description"
		infoname := "[기생충] 30초 예고.info.json"
		copyFileContents(filepath.Join(config.RootDir, "testdata", descname), filepath.Join(dir, descname))
		copyFileContents(filepath.Join(config.RootDir, "testdata", infoname), filepath.Join(dir, infoname))
		return &Cmd{Content: []byte{}}
	} else if name == "youtube-dl" &&
		arg[0] == "--skip-download" &&
		arg[1] == "--output" &&
		arg[2][:4] == "/var" && // /var/folders/53/tpn8zp511y1gdz9k_srhvbh80000gn/T/317835077/%(title)s.%(ext)s
		arg[3] == "--write-description" &&
		arg[4] == "--write-info-json" &&
		arg[5] == "--write-annotations" &&
		arg[6] == "--write-sub" &&
		arg[7] == "--" &&
		arg[8] == "-BIDXOp6_LA" {
		dir := filepath.Dir(arg[2])
		descname := "Go Modules - Dependency Management the Right Way.description"
		infoname := "Go Modules - Dependency Management the Right Way.info.json"
		err := copyFileContents(filepath.Join(config.RootDir, "testdata", descname), filepath.Join(dir, descname))
		if err != nil {
			log.Fatalf("could not copy description: %v", err)
		}
		err = copyFileContents(filepath.Join(config.RootDir, "testdata", infoname), filepath.Join(dir, infoname))
		if err != nil {
			log.Fatalf("could not copy info: %v", err)
		}
		return &Cmd{Content: []byte{}}
	} else if name == "youtube-dl" &&
		arg[0] == "--format" &&
		arg[1] == "18" &&
		arg[2] == "--output" &&
		arg[3][:4] == "/var" && // /var/folders/53/tpn8zp511y1gdz9k_srhvbh80000gn/T/407901467/%(title)s_%(format_note)s.%(ext)s
		arg[4] == "--get-filename" &&
		arg[5] == "--" &&
		arg[6] == "GSVsfCCtRr0" {
		dir := filepath.Dir(arg[3])
		return &Cmd{Content: []byte(filepath.Join(dir, "[기생충] 30초 예고_360p.mp4"))}
	} else if name == "youtube-dl" &&
		arg[0] == "--format" &&
		arg[1] == "18" &&
		arg[2] == "--output" &&
		arg[3][:4] == "/var" && // /var/folders/53/tpn8zp511y1gdz9k_srhvbh80000gn/T/407901467/%(title)s_%(format_note)s.%(ext)s
		arg[4] == "--" &&
		arg[5] == "GSVsfCCtRr0" {
		dir := filepath.Dir(arg[3])
		videoname := "[기생충] 30초 예고_360p.mp4"
		err := copyFileContents(filepath.Join(config.RootDir, "testdata", videoname), filepath.Join(dir, videoname))
		if err != nil {
			log.Fatalf("could not copy video: %s", err)
		}
		return &Cmd{Content: []byte{}}
	} else {
		log.Fatalf("could not mock command %s %s", name, arg)
		return nil
	}
}

type Cmd struct{ Content []byte }

func (c *Cmd) Output() ([]byte, error)         { return c.Content, nil }
func (c *Cmd) CombinedOutput() ([]byte, error) { return c.Content, nil }

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
