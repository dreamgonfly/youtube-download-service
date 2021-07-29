package storage

type StorerMock struct{}

func (s *StorerMock) DownloadFile(uri string) ([]byte, error) {
	return nil, nil
}

func (s *StorerMock) UploadFile(path string, uri string) error {
	return nil
}
