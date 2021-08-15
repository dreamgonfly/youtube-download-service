package videoinfo

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type Info struct {
	Title       string
	DurationSecond float64
	Formats     []Format
	Thumbnails  []Thumbnail
}

type Format struct {
	Filesize   int64
	FormatId   string
	FormatNote string
	Ext        string
}

type Thumbnail struct {
	Height     int64
	Width      int64
	URL        string
	Resolution string
	Id         int64
}

func NewInfo(path string) (Info, error) {
	infoFile, err := os.Open(path)
	if err != nil {
		return Info{}, errors.Wrap(err, "could not open info file")
	}
	defer infoFile.Close()

	infoBytes, err := ioutil.ReadAll(infoFile)
	if err != nil {
		return Info{}, errors.Wrap(err, "could not read info")
	}

	infoJson := string(infoBytes)

	title := gjson.Get(infoJson, "title").String()
	duration := gjson.Get(infoJson, "duration").Float()

	var formats []Format
	formatsResult := gjson.Get(infoJson, "formats")
	formatsResult.ForEach(func(key, value gjson.Result) bool {
		format := Format{
			Filesize:   value.Get("filesize").Int(),
			FormatId:   value.Get("format_id").String(),
			FormatNote: value.Get("format_note").String(),
			Ext:        value.Get("ext").String(),
		}
		formats = append(formats, format)
		return true // keep iterating
	})

	if len(formats) == 0 {
		return Info{}, errors.New("formats are empty")
	}

	var thumbnails []Thumbnail
	thumbnailResult := gjson.Get(infoJson, "thumbnails")
	thumbnailResult.ForEach(func(key, value gjson.Result) bool {
		thumbnail := Thumbnail{
			Height:     value.Get("height").Int(),
			Width:      value.Get("width").Int(),
			URL:        value.Get("url").String(),
			Resolution: value.Get("resolution").String(),
			Id:         value.Get("id").Int(),
		}
		thumbnails = append(thumbnails, thumbnail)
		return true // keep iterating
	})

	if len(thumbnails) == 0 {
		return Info{}, errors.New("thumbnails are empty")
	}

	info := Info{
		Title:       title,
		DurationSecond: duration,
		Formats:     formats,
		Thumbnails:  thumbnails,
	}
	return info, nil
}

func EstimateFilesize(formatNote string, durationSec float64) (int64, error) {
	const BytesPerSecond720p = 68649
	const BytesPerSecond360p = 58993
	if formatNote == "360p" {
		return int64(durationSec * BytesPerSecond360p), nil
	} else if formatNote == "720p" {
		return int64(durationSec * BytesPerSecond720p), nil
	} else {
		return 0, errors.New("unknown format_note")
	}
}
