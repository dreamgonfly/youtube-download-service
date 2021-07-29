package videodownload

type Downloader interface {
	Preview(id string, dir string) (description string, info string, thumbnail string, err error)
	Download(id string, format_code string, dir string) (video string, err error)
}
