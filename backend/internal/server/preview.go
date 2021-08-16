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
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		descriptionPath, infoPath, err := s.youtubedl.Preview(id, tempDir)
		if err != nil {
			err = errors.Wrap(err, "could not download preview")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uploadDescriptionDone := make(chan struct{})
		uploadInfoDone := make(chan struct{})
		go s.UploadLocalFile(id, descriptionPath, uploadDescriptionDone)
		go s.UploadLocalFile(id, infoPath, uploadInfoDone)
		defer s.CleanupTempDir(tempDir, uploadDescriptionDone, uploadInfoDone)

		info, err := videoinfo.NewInfo(infoPath)
		if err != nil {
			err = errors.Wrap(err, "could not extract info")
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
				if err != nil {
					err = errors.Wrap(err, "could not estimate file size")
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			estimatedFormats = append(estimatedFormats, format)
		}

		lastThumbnail := info.Thumbnails[len(info.Thumbnails)-1]

		output := struct {
			Title          string
			DurationSecond float64
			Thumbnail      string
			Name           string
			Formats        []videoinfo.Format
		}{
			Title:          info.Title,
			DurationSecond: info.DurationSecond,
			Thumbnail:      lastThumbnail.URL,
			Name:           youtubefile.Stem(filepath.Base(descriptionPath)),
			Formats:        estimatedFormats,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) CleanupTempDir(tempDir string, uploadDescriptionDone chan struct{}, uploadInfoDone chan struct{}) {
	<-uploadDescriptionDone
	<-uploadInfoDone
	err := os.RemoveAll(tempDir)
	if err != nil {
		err = errors.Wrap(err, "could not remove tempDir")
		log.Println(err)
	}
}

func (s *Server) UploadLocalFile(id string, path string, done chan struct{}) {
	key := filepath.Join("videos", id, filepath.Base(path))
	err := gcs.UploadFile(s.context, s.gcsClient, path, key)
	if err != nil {
		err = errors.Wrap(err, "could not upload to gcs")
		log.Println(err)
	}
	close(done)
}
