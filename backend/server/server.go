package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"youtube-download-backend/extract"
	"youtube-download-backend/storage"
	"youtube-download-backend/videodownload"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	router     *mux.Router
	storage    storage.Storer
	downloader videodownload.Downloader
}

func NewServer(st storage.Storer, dw videodownload.Downloader) *Server {
	router := mux.NewRouter()
	s := &Server{
		router:     router,
		storage:    st,
		downloader: dw,
	}
	s.routes()
	return s
}

// Make Server an http.Handle
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!\n")
	}
}

func (s *Server) handlePreview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			err := errors.New("video id not provided")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		mediatype, _, err := mime.ParseMediaType(r.Header.Get("Accept"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}
		if mediatype != "multipart/form-data" {
			http.Error(w, "set Accept: multipart/form-data", http.StatusMultipleChoices)
			return
		}

		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			err = errors.Wrap(err, "could not create temp dir")
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.RemoveAll(tempDir)

		description, info, thumbnail, err := s.downloader.Preview(id, tempDir)
		if err != nil {
			err = errors.Wrap(err, "could not download preview")
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		formats, err := extract.ExtractFormatsFromInfo(info)
		if err != nil {
			err = errors.Wrap(err, "could not extract formats")
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

		mw := multipart.NewWriter(w)
		w.Header().Set("Content-Type", mw.FormDataContentType())
		fw, err := mw.CreateFormField("formats")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(fw).Encode(estimatedFormats)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bytes, err := ioutil.ReadFile(thumbnail)
		fw, err = mw.CreateFormFile("thumbnail", path.Base(thumbnail))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = fw.Write(bytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := mw.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		key := path.Join("videos", id, path.Base(description))
		uri := storage.ComposeURI("gs", "youtube-download-backend-beta", key)
		err = s.storage.UploadFile(description, uri)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		key = path.Join("videos", id, path.Base(info))
		uri = storage.ComposeURI("gs", "youtube-download-backend-beta", key)
		err = s.storage.UploadFile(info, uri)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		key = path.Join("videos", id, path.Base(thumbnail))
		uri = storage.ComposeURI("gs", "youtube-download-backend-beta", key)
		err = s.storage.UploadFile(info, uri)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (s *Server) handleDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			err := errors.New("video id is missing")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		formatCodes, ok := r.URL.Query()["format_code"]
		if !ok || len(formatCodes[0]) < 1 {
			http.Error(w, "format code is missing", http.StatusBadRequest)
			return
		}
		formatCode := formatCodes[0]

		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			err = errors.Wrap(err, "could not create temp dir")
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.RemoveAll(tempDir)

		video, err := s.downloader.Download(id, formatCode, tempDir)
		if err != nil {
			err = errors.Wrap(err, "could not download video")
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		key := path.Join("videos", id, path.Base(video))
		uri := storage.ComposeURI("gs", "youtube-download-backend-beta", key)
		err = s.storage.UploadFile(video, uri)
		if err != nil {
			err = errors.Wrap(err, "could not upload video")
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filepath.Base(video)))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, video)
	}
}
