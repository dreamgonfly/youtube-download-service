package extract

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

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

func ExtractFormatsFromInfo(path string) ([]Format, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open info file")
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not read info")
	}

	var formats []Format
	result := gjson.Get(string(byteValue), "formats")
	result.ForEach(func(key, value gjson.Result) bool {
		format := Format{
			Filesize:   value.Get("filesize").Int(),
			FormatId:   value.Get("format_id").String(),
			FormatNote: value.Get("format_note").String(),
			Ext:        value.Get("ext").String(),
		}
		formats = append(formats, format)
		return true // keep iterating
	})
	return formats, nil
}

func ExtractDuration(info string) (float64, error) {
	jsonFile, err := os.Open(info)
	if err != nil {
		return 0, errors.Wrap(err, "could not open info file")
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return 0, errors.Wrap(err, "could not read info")
	}

	result := gjson.Get(string(byteValue), "duration").Float()
	return result, nil
}

func ExtractThumbnails(info string) ([]Thumbnail, error) {
	jsonFile, err := os.Open(info)
	if err != nil {
		return nil, errors.Wrap(err, "could not open info file")
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not read info")
	}

	var thumbnails []Thumbnail
	result := gjson.Get(string(byteValue), "thumbnails")
	result.ForEach(func(key, value gjson.Result) bool {
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

	return thumbnails, nil
}

func EstimateFilesize(formatNote string, info string) (int64, error) {
	const BytesPerSecond720p = 68649
	const BytesPerSecond360p = 58993
	duration, err := ExtractDuration(info)
	if err != nil {
		return 0, errors.Wrap(err, "could not extract duration")
	}
	if formatNote == "360p" {
		return int64(duration * BytesPerSecond360p), nil
	} else if formatNote == "720p" {
		return int64(duration * BytesPerSecond720p), nil
	} else {
		return 0, errors.New("unknown format_note")
	}
}
