package gcs

import (
	"errors"
	"fmt"

	"cloud.google.com/go/storage"
)

func SignedURL(bucket, name string, opts *storage.SignedURLOptions) (string, error) {
	if name == "videos/GSVsfCCtRr0/[기생충] 30초 예고.jpg" {
		return "https://storage.googleapis.com/youtube-download-backend-beta/videos/GSVsfCCtRr0/%5B%EA%B8%B0%EC%83%9D%EC%B6%A9%5D%2030%EC%B4%88%20%EC%98%88%EA%B3%A0.info.jpg?X-Goog-Algorithm=GOOG4-RSA-SHA256\u0026X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210801%2Fauto%2Fstorage%2Fgoog4_request\u0026X-Goog-Date=20210801T084013Z\u0026X-Goog-Expires=899\u0026X-Goog-Signature=31390355ce8353e169279bc4286078c29ef2ed5bcb2bcc13380c67efac26e4155275e0d1a3683cc77e7d9188e33ffaee0e7056c584951e339c327934e7d08a78a0e549acd05b1506ed62c6e048c8500416cc7086f07bd9a03fa1df6b678220b0a1c810d0715f904c889a80fc3208a997ebb6e117f6857cf526daa32bb515c048384f85b27456c6c39ae7647efe2f0be24f0905a52fe0e5b94bbd6871a71c34951f1317cbccee23fa4ed0566415e9b55221fba05558ebdbc9c6e48f11396e1766f7f62b49756b8daa7cc54c6169176f81e6726e91e02dec30d671d7cb4cca85cc6556c5ef312293ece4b22792b53aa05481e3ad0ba225ca3435e935adb25a7143\u0026X-Goog-SignedHeaders=host", nil
	} else {
		return "", errors.New(fmt.Sprintf("could not mock %s", name))
	}
}
