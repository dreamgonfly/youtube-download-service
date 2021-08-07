package exec

import (
	"bytes"
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
		arg[0] == "--format" &&
		arg[1] == "18" &&
		arg[2] == "--output" &&
		arg[3] == "%(title)s_%(format_note)s.%(ext)s" &&
		arg[4] == "--get-filename" &&
		arg[5] == "--" &&
		arg[6] == "GSVsfCCtRr0" {
		return &Cmd{Content: []byte("[기생충] 30초 예고_360p.mp4")}
	} else if name == "youtube-dl" &&
		arg[0] == "--format" &&
		arg[1] == "18" &&
		arg[2] == "--get-url" &&
		arg[3] == "--" &&
		arg[4] == "GSVsfCCtRr0" {
		return &Cmd{Content: []byte("https://r5---sn-3u-bh2ls.googlevideo.com/videoplayback?expire=1628365436&ei=G44OYYn3OpnK2roPqNG1wAo&ip=220.85.16.207&id=o-APCYE7oLizYKGaPmPdASOA59Dnlrr3DK3AO6TD5Rjbzq&itag=18&source=youtube&requiressl=yes&mh=l0&mm=31%2C26&mn=sn-3u-bh2ls%2Csn-npoe7ns7&ms=au%2Conr&mv=m&mvi=5&pl=19&initcwndbps=1170000&vprv=1&mime=video%2Fmp4&ns=pVJ9QNUMaUC_UYjGA8J-F2EG&gir=yes&clen=1348634&ratebypass=yes&dur=30.139&lmt=1557985846949015&mt=1628343659&fvip=5&fexp=24001373%2C24007246&c=WEB&txp=2211222&n=Uygq1sf6iX2k6tw&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cvprv%2Cmime%2Cns%2Cgir%2Cclen%2Cratebypass%2Cdur%2Clmt&sig=AOq0QJ8wRAIgX9wE5bvmulhalbftGUnNh2DU9NuzlL16VTJIpEKwCUsCIGUCRhLfnUY5NvOHQR6CokGp_4bKmQ1lz4a2sm4Pb2wJ&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=AG3C_xAwRAIgfCmTn87LXf9kALwVgWwGMc3Zj3MkwcQumr-8-bzaPvMCICs-pF19FYYB2PmNDK53M7NevXNQdmBoLthwUbZtRUzo")}
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
		arg[2] == "--force-ipv4" &&
		arg[3] == "--output" &&
		arg[4] == "-" &&
		arg[5] == "--" &&
		arg[6] == "GSVsfCCtRr0" {
		videoname := "[기생충] 30초 예고_360p.mp4"
		videopath := filepath.Join(config.RootDir, "testdata", videoname)
		in, err := os.Open(videopath)
		if err != nil {
			log.Fatalf("could not open video: %s", err)
		}
		content, err := io.ReadAll(in)
		if err != nil {
			log.Fatalf("could not read video: %s", err)
		}
		return &Cmd{Content: content}
	} else {
		log.Fatalf("could not mock command %s %s", name, arg)
		return nil
	}
}

type Cmd struct{ Content []byte }

func (c *Cmd) Output() ([]byte, error)         { return c.Content, nil }
func (c *Cmd) CombinedOutput() ([]byte, error) { return c.Content, nil }
func (c *Cmd) StdoutPipe() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(c.Content)), nil
}
func (c *Cmd) StderrPipe() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte{})), nil
}
func (c *Cmd) Start() error   { return nil }
func (c *Cmd) Wait() error    { return nil }
func (c *Cmd) String() string { return "" }

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
