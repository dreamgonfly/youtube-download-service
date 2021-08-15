package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"youtube-download-backend/internal/gcs"
	"youtube-download-backend/internal/videoinfo"
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

		descriptionPath, infoPath, err := s.youtubedl.Preview(id, tempDir)
		if err != nil {
			err = errors.Wrap(err, "could not download preview")
			// TODO: logging
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go func() {
			key := filepath.Join("videos", id, filepath.Base(descriptionPath))
			err = gcs.UploadFile(s.context, s.gcsClient, descriptionPath, key)
			if err != nil {
				err = errors.Wrap(err, "could not upload description")
				log.Println(err)
				// TODO: logging
				return
			}
			close(uploadDescriptionDone)
		}()
		go func() {
			key := filepath.Join("videos", id, filepath.Base(infoPath))
			err = gcs.UploadFile(s.context, s.gcsClient, infoPath, key)
			if err != nil {
				err = errors.Wrap(err, "could not upload info")
				log.Println(err)
				// TODO: logging
				return
			}
			close(uploadInfoDone)
		}()

		info, err := videoinfo.NewInfo(infoPath)
		if err != nil {
			err = errors.Wrap(err, "could not extract info")
			// TODO: logging
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var filteredFormats []videoinfo.Format
		for _, format := range info.Formats {
			if format.FormatId == "18" || format.FormatId == "22" {
				filteredFormats = append(filteredFormats, format)
			}
		}

		var estimatedFormats []videoinfo.Format
		for _, format := range filteredFormats {
			if format.Filesize == 0 {
				format.Filesize, err = videoinfo.EstimateFilesize(format.FormatNote, info.DurationSecond)
				// TODO: log estimation err
			}
			estimatedFormats = append(estimatedFormats, format)
		}

		lastThumbnail := info.Thumbnails[len(info.Thumbnails)-1]

		type Output struct {
			Title          string
			DurationSecond float64
			Thumbnail      string
			Name           string
			Formats        []videoinfo.Format
		}

		o := Output{Title: info.Title, DurationSecond: info.DurationSecond, Thumbnail: lastThumbnail.URL, Name: youtubefile.Stem(filepath.Base(descriptionPath)), Formats: estimatedFormats}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(o)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
