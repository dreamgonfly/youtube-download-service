package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"youtube-download-backend/internal/extract"
	"youtube-download-backend/internal/gcs"
	"youtube-download-backend/internal/youtubefile"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *Server) handlePreview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			err := errors.New("video id not provided")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			err = errors.Wrap(err, "could not create temp dir")
			// TODO: logging
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		uploadDescriptionDone := make(chan struct{})
		uploadInfoDone := make(chan struct{})
		defer func() {
			<-uploadDescriptionDone
			<-uploadInfoDone
			err := os.RemoveAll(tempDir)
			if err != nil {
				err = errors.Wrap(err, "could not remove tempDir")
				log.Println(err)
			}
		}()

		description, info, err := s.youtubedl.Preview(id, tempDir)
		if err != nil {
			err = errors.Wrap(err, "could not download preview")
			// TODO: logging
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		formats, err := extract.ExtractFormatsFromInfo(info)
		if err != nil {
			err = errors.Wrap(err, "could not extract formats")
			// TODO: logging
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		thumbnails, err := extract.ExtractThumbnails(info)
		if err != nil {
			err = errors.Wrap(err, "could not extract thumbnails")
			// TODO: logging
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		last_thumbnail := thumbnails[len(thumbnails)-1]
		var filteredFormats []extract.Format
		for _, format := range formats {
			if format.FormatId == "18" || format.FormatId == "22" {
				filteredFormats = append(filteredFormats, format)
			}
		}

		var estimatedFormats []extract.Format
		for _, format := range filteredFormats {
			if format.Filesize == 0 {
				format.Filesize, err = extract.EstimateFilesize(format.FormatNote, info)
				// TODO: log estimation err
			}
			estimatedFormats = append(estimatedFormats, format)
		}

		go func() {
			key := filepath.Join("videos", id, filepath.Base(description))
			err = gcs.UploadFile(s.context, s.gcsClient, description, key)
			if err != nil {
				err = errors.Wrap(err, "could not upload description")
				log.Println(err)
				// TODO: logging
				return
			}
			close(uploadDescriptionDone)
		}()
		go func() {
			key := filepath.Join("videos", id, filepath.Base(info))
			err = gcs.UploadFile(s.context, s.gcsClient, info, key)
			if err != nil {
				err = errors.Wrap(err, "could not upload info")
				log.Println(err)
				// TODO: logging
				return
			}
			close(uploadInfoDone)
		}()

		type Output struct {
			Thumbnail string
			Name      string
			Formats   []extract.Format
		}
		o := Output{Thumbnail: last_thumbnail.URL, Name: youtubefile.Stem(filepath.Base(description)), Formats: estimatedFormats}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(o)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
