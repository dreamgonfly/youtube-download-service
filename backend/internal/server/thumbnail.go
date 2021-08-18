package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
	"youtube-download-backend/internal/gcs"
	"youtube-download-backend/internal/httpfile"
	"youtube-download-backend/internal/logging"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type ThumbnailRequest struct {
	URL  string
	Name string
}

// handleUpdateThumbnail stores thumbnail in GCS then returns signed url
func (s *Server) handleUpdateThumbnail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			err := errors.New("video id is missing")
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		var thumbnailRequest ThumbnailRequest
		err := json.NewDecoder(r.Body).Decode(&thumbnailRequest)
		defer r.Body.Close()
		if err != nil {
			err := errors.Wrap(err, "could not parse request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tempFile, err := ioutil.TempFile("", "")
		if err != nil {
			err = errors.Wrap(err, "could not create temp file")
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		u, err := url.Parse(thumbnailRequest.URL)
		if err != nil {
			err = errors.Wrapf(err, "could not parse thumbnail url (%s)", thumbnailRequest.URL)
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = httpfile.DownloadFile(s.httpClient, thumbnailRequest.URL, tempFile.Name())
		if err != nil {
			err = errors.Wrapf(err, "could not download thumbnail url (%s)", thumbnailRequest.URL)
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		thumbnailKey := filepath.Join("videos", id, thumbnailRequest.Name+filepath.Ext(u.Path))
		err = gcs.UploadFile(s.context, s.gcsClient, tempFile.Name(), thumbnailKey)
		if err != nil {
			err = errors.Wrap(err, "could not upload thumbnail")
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		signedURL, err := gcs.GenerateV4GetObjectSignedURL(s.signFunc, thumbnailKey, 15*time.Minute)
		if err != nil {
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		output := struct {
			Thumbnail string
		}{
			Thumbnail: signedURL,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tempFile.Close(); err != nil {
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
