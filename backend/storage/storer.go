package storage

type Storer interface {
	DownloadFile(uri string) ([]byte, error)
	UploadFile(path string, uri string) error
}
