package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"youtube-download-backend/internal/youtubefile"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const BUFFER_SIZE = 1024

func (s *Server) handleDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			err := errors.New("video id is missing")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		formats, ok := r.URL.Query()["format"]
		if !ok || len(formats[0]) < 1 {
			log.Println("format is missing")
			http.Error(w, "format is missing", http.StatusBadRequest)
			return
		}
		format := formats[0]

		name, err := s.youtubedl.GetNameWithFormat(id, format)
		if err != nil {
			err = errors.Wrap(err, "could not get name")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cmd := s.youtubedl.StreamDownloadCommand(id, format, w)
		err = s.StreamDownload(cmd, name, w)
		if err != nil {
			err = errors.Wrap(err, "could not download video")
			log.Println(err)
			// TODO: logging
			w.Header().Del("Content-Disposition")
			w.Header().Del("Content-Type")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) StreamDownload(cmd youtubefile.Outputer, name string, w http.ResponseWriter) error {
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(name))
	w.Header().Set("Content-Type", "application/octet-stream")

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
		w.Write(data)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		} else {
			return errors.New("could not flush http")
		}
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
	return nil
}
