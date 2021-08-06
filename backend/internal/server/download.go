package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *Server) handleDownload() http.HandlerFunc {
	// TODo: get name from file
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			err := errors.New("video id is missing")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		formatCodes, ok := r.URL.Query()["format"]
		if !ok || len(formatCodes[0]) < 1 {
			log.Println("format is missing")
			http.Error(w, "format is missing", http.StatusBadRequest)
			return
		}
		formatCode := formatCodes[0]

		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			err = errors.Wrap(err, "could not create temp dir")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		uploadDone := make(chan struct{})
		go func() {
			<-uploadDone
			err := os.RemoveAll(tempDir)
			if err != nil {
				err = errors.Wrap(err, "could not remove tempDir")
				log.Println(err)
			}
		}()

		err = s.youtubedl.DownloadStream(id, formatCode, w)
		if err != nil {
			err = errors.Wrap(err, "could not download video")
			log.Println(err)
			// TODO: logging
			w.Header().Del("Content-Disposition")
			w.Header().Del("Content-Type")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: streaming to GCS
		// go func() {
		// 	key := filepath.Join("videos", id, filepath.Base(video))
		// 	err = gcs.UploadFile(s.context, s.gcsClient, video, key)
		// 	if err != nil {
		// 		err = errors.Wrap(err, "could not upload video")
		// 		log.Println(err)
		// 		// TODO: logging
		// 	}
		// 	close(uploadDone)
		// }()
	}
}
