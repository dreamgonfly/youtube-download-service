package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	. "youtube-download-backend/internal/config"
	"youtube-download-backend/internal/extract"
	"youtube-download-backend/internal/gcs"
	"youtube-download-backend/internal/httpfile"
	"youtube-download-backend/internal/storagepath"
	"youtube-download-backend/internal/youtubefile"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

type Server struct {
	router     *mux.Router
	context    context.Context
	youtubedl  youtubefile.YoutubeDl
	gcsClient  gcs.Clienter
	httpClient httpfile.Clienter
}

func NewServer(ctx context.Context, c youtubefile.Commander, g gcs.Clienter, h httpfile.Clienter) *Server {
	r := mux.NewRouter()
	s := &Server{
		router:     r,
		context:    ctx,
		youtubedl:  youtubefile.YoutubeDl{ExecCommand: c},
		gcsClient:  g,
		httpClient: h,
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
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		mediatype, _, err := mime.ParseMediaType(r.Header.Get("Accept"))
		if err != nil {
			log.Println(err)
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
		thumbnailPath := filepath.Join(tempDir, strings.Join([]string{youtubefile.Stem(filepath.Base(info)), filepath.Ext(u.Path)}, ""))
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

		mw := multipart.NewWriter(w)
		w.Header().Set("Content-Type", mw.FormDataContentType())
		fw, err := mw.CreateFormField("formats")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(fw).Encode(estimatedFormats)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bytes, err := ioutil.ReadFile(thumbnailPath)
		fw, err = mw.CreateFormFile("thumbnail", filepath.Base(thumbnailPath))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = fw.Write(bytes)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := mw.Close(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		key := filepath.Join("videos", id, slug.MakeLang(youtubefile.Stem(filepath.Base(description)), "en")+filepath.Ext(description))
		uri := storagepath.ComposeURI("gs", Conf.Bucket, key)
		err = gcs.UploadFile(s.context, s.gcsClient, description, uri)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		key = filepath.Join("videos", id, slug.MakeLang(youtubefile.Stem(filepath.Base(info)), "en")+filepath.Ext(info))
		uri = storagepath.ComposeURI("gs", Conf.Bucket, key)
		err = gcs.UploadFile(s.context, s.gcsClient, info, uri)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			log.Println(err)
			// TODO: logging
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		key = filepath.Join("videos", id, slug.MakeLang(youtubefile.Stem(filepath.Base(thumbnailPath)), "en")+filepath.Ext(thumbnailPath))
		uri = storagepath.ComposeURI("gs", Conf.Bucket, key)
		err = gcs.UploadFile(s.context, s.gcsClient, thumbnailPath, uri)
		if err != nil {
			err = errors.Wrap(err, "could not upload preview")
			log.Println(err)
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
		key := filepath.Join("videos", id, slug.MakeLang(youtubefile.Stem(filepath.Base(video)), "en")+filepath.Ext(video))
		uri := storagepath.ComposeURI("gs", Conf.Bucket, key)
		err = gcs.UploadFile(s.context, s.gcsClient, video, uri)
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
