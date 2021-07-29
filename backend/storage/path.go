package storage

import (
	"errors"
	"net/url"
	"strings"
)

func ParseURI(uri string) (scheme, bucket, key string, err error) {
	// TODO: Deal with fragments
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", "", err
	}
	if u.Path == "" {
		return u.Scheme, u.Host, "", errors.New("not a valid uri")
	}
	return u.Scheme, u.Host, strings.TrimPrefix(u.Path, "/"), nil
}

func ComposeURI(scheme, bucket, key string) (uri string) {
	return scheme + "://" + bucket + "/" + key
}
