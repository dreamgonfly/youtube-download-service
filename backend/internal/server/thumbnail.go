package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"youtube-download-backend/internal/gcs"
	"youtube-download-backend/internal/httpfile"

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
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		var t ThumbnailRequest
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			err := errors.Wrap(err, "could not parse request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tf, err := ioutil.TempFile("", "")
		if err != nil {
			err = errors.Wrap(err, "could not create temp dir")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.Remove(tf.Name())
		u, err := url.Parse(t.URL)
		if err != nil {
			err = errors.Wrapf(err, "could not parse thumbnail url (%s)", t.URL)
			// TODO: logging
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = httpfile.DownloadFile(s.httpClient, t.URL, tf.Name())
		if err != nil {
			err = errors.Wrapf(err, "could not download thumbnail url (%s)", t.URL)
			// TODO: logging
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		thumbnailKey := filepath.Join("videos", id, t.Name+filepath.Ext(u.Path))
		err = gcs.UploadFile(s.context, s.gcsClient, tf.Name(), thumbnailKey)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		type Output struct {
			Thumbnail string
		}

		signedURL, err := gcs.GenerateV4GetObjectSignedURL(s.signFunc, thumbnailKey)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		o := Output{Thumbnail: signedURL}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(o)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tf.Close(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
