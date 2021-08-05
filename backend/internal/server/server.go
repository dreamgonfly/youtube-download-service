package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"youtube-download-backend/internal/extract"
	"youtube-download-backend/internal/gcs"
	"youtube-download-backend/internal/httpfile"
	"youtube-download-backend/internal/youtubefile"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	router     *mux.Router
	context    context.Context
	youtubedl  youtubefile.YoutubeDl
	gcsClient  gcs.Clienter
	httpClient httpfile.Clienter
	signFunc   gcs.SignFunc
}

func NewServer(ctx context.Context, c youtubefile.Commander, g gcs.Clienter, h httpfile.Clienter, sf gcs.SignFunc) *Server {
	r := mux.NewRouter()
	s := &Server{
		router:     r,
		context:    ctx,
		youtubedl:  youtubefile.YoutubeDl{ExecCommand: c},
		gcsClient:  g,
		httpClient: h,
		signFunc:   sf,
	}
	s.routes()
	return s
}

// Make Server an http.Handle
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(&w)
	if r.Method == "OPTIONS" { // Handling pre-flight OPTIONS requests
		return
	}
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
		defer os.RemoveAll(tempDir)

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
		u, err := url.Parse(last_thumbnail.URL)
		if err != nil {
			err = errors.Wrapf(err, "could not parse thumbnail url (%s)", last_thumbnail.URL)
			// TODO: logging
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		thumbnailPath := filepath.Join(tempDir, strings.Join([]string{youtubefile.Stem(filepath.Base(description)), filepath.Ext(u.Path)}, ""))
		err = httpfile.DownloadFile(s.httpClient, last_thumbnail.URL, thumbnailPath)
		if err != nil {
			err = errors.Wrapf(err, "could not download thumbnail url (%s)", last_thumbnail.URL)
			// TODO: logging
			log.Println(err)
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

		key := filepath.Join("videos", id, filepath.Base(description))
		err = gcs.UploadFile(s.context, s.gcsClient, description, key)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		key = filepath.Join("videos", id, filepath.Base(info))
		err = gcs.UploadFile(s.context, s.gcsClient, info, key)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		thumbnailKey := filepath.Join("videos", id, filepath.Base(thumbnailPath))
		err = gcs.UploadFile(s.context, s.gcsClient, thumbnailPath, thumbnailKey)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type Output struct {
			Thumbnail string
			Formats   []extract.Format
		}

		signedURL, err := gcs.GenerateV4GetObjectSignedURL(s.signFunc, thumbnailKey)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		o := Output{Thumbnail: signedURL, Formats: estimatedFormats}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(o)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

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
		defer os.RemoveAll(tempDir)

		video, err := s.youtubedl.Download(id, formatCode, tempDir)
		if err != nil {
			err = errors.Wrap(err, "could not download video")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		key := filepath.Join("videos", id, filepath.Base(video))
		err = gcs.UploadFile(s.context, s.gcsClient, video, key)
		if err != nil {
			err = errors.Wrap(err, "could not upload video")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filepath.Base(video)))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, video)
	}
}
