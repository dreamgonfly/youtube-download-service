package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	. "youtube-download-backend/internal/config"
	"youtube-download-backend/internal/gcs"
	"youtube-download-backend/internal/logging"
	"youtube-download-backend/internal/youtubefile"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *Server) handlePlay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			err := errors.New("video id is missing")
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		formats, ok := r.URL.Query()["format"]
		if !ok || len(formats) != 1 {
			err := errors.New("format is missing")
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		format := formats[0]

		var filename string
		filenames, ok := r.URL.Query()["filename"]
		if ok && len(filenames) == 1 {
			filename = filenames[0]
		} else {
			name, err := s.youtubedl.GetFilenameWithFormat(id, format)
			if err != nil {
				err = errors.Wrap(err, "could not get name")
				logging.Logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			filename = name
		}

		cmd := s.youtubedl.StreamDownloadCommand(id, format, w)

		key := filepath.Join("videos", id, filename)
		err := s.StreamSave(cmd, key)
		if err != nil {
			err = errors.Wrap(err, "could not save video")
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		signedURL, err := gcs.GenerateV4GetObjectSignedURL(s.signFunc, key, 3*time.Hour)

		output := struct {
			Filename string
			URL      string
		}{
			Filename: filename,
			URL:      signedURL,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			logging.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (s *Server) StreamSave(cmd youtubefile.Outputer, key string) error {
	ctx, cancel := context.WithTimeout(s.context, 1*time.Hour)
	defer cancel()

	wc := s.gcsClient.Bucket(Config.Bucket).Object(key).NewWriter(ctx)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("stdout error command (%s)", cmd.String()))
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("stderr error command (%s)", cmd.String()))
	}

	err = cmd.Start()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("error starting (%s)", cmd.String()))
		return err
	}
	buffer := make([]byte, BUFFER_SIZE)
	for {
		n, err := stdout.Read(buffer)
		if err != nil {
			stdout.Close()
			break
		}
		data := buffer[0:n]
		wc.Write(data)

		// reset buffer
		for i := 0; i < n; i++ {
			buffer[i] = 0
		}
	}
	errout, err := io.ReadAll(stderr)
	if err != nil {
		err = errors.Wrap(err, "could not read stderr")
	}
	err = cmd.Wait()
	if err != nil {
		err = errors.Wrap(err, strings.TrimSpace(string(errout)))
		return errors.Wrap(err, fmt.Sprintf("error waiting command (%s)", cmd.String()))
	}

	err = wc.Close()
	if err != nil {
		return errors.Wrap(err, "Writer.Close")
	}
	return nil
}
