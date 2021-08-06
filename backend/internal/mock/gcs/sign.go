package gcs

import (
	"errors"
	"fmt"

	"cloud.google.com/go/storage"
)

func SignedURL(bucket, name string, opts *storage.SignedURLOptions) (string, error) {
	if name == "videos/GSVsfCCtRr0/[기생충] 30초 예고.jpg" {
		return "https://storage.googleapis.com/youtube-download-backend-beta/videos/GSVsfCCtRr0/%5B%EA%B8%B0%EC%83%9D%EC%B6%A9%5D%2030%EC%B4%88%20%EC%98%88%EA%B3%A0.jpg?X-Goog-Algorithm=GOOG4-RSA-SHA256\u0026X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210801%2Fauto%2Fstorage%2Fgoog4_request\u0026X-Goog-Date=20210801T084013Z\u0026X-Goog-Expires=899\u0026X-Goog-Signature=31390355ce8353e169279bc4286078c29ef2ed5bcb2bcc13380c67efac26e4155275e0d1a3683cc77e7d9188e33ffaee0e7056c584951e339c327934e7d08a78a0e549acd05b1506ed62c6e048c8500416cc7086f07bd9a03fa1df6b678220b0a1c810d0715f904c889a80fc3208a997ebb6e117f6857cf526daa32bb515c048384f85b27456c6c39ae7647efe2f0be24f0905a52fe0e5b94bbd6871a71c34951f1317cbccee23fa4ed0566415e9b55221fba05558ebdbc9c6e48f11396e1766f7f62b49756b8daa7cc54c6169176f81e6726e91e02dec30d671d7cb4cca85cc6556c5ef312293ece4b22792b53aa05481e3ad0ba225ca3435e935adb25a7143\u0026X-Goog-SignedHeaders=host", nil
	} else if name == "videos/-BIDXOp6_LA/Go Modules - Dependency Management the Right Way.webp" {
		return "https://storage.googleapis.com/youtube-download-backend-beta/videos/-BIDXOp6_LA/Go%20Modules%20-%20Dependency%20Management%20the%20Right%20Way.webp?X-Goog-Algorithm=GOOG4-RSA-SHA256\u0026X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210801%2Fauto%2Fstorage%2Fgoog4_request\u0026X-Goog-Date=20210801T201751Z\u0026X-Goog-Expires=899\u0026X-Goog-Signature=6baaea933a08dc902f21a4b95f5d88b7fa42664d0ffd83385a24b180c248ab5751c8cc9a9f927517f29235751fef7606825210340e7f995fb4d1609a9c78d153062cb7fa11f67e082c28f50f6262632c3337bc225584b4c15405acfb4b6a03e4d253db41b14d39113bce36140c4afae634a8a9e51dfd08f54700c1512996857dfe6604ecb335228e4baecce9458b160537f97d5f2900448a1edf2da3d2c57da2db7690b8d8c7108762cdf4123ee4cb718352859a0181879bb7cd38ba4b5de679a7fa79ad0ac097819af5910dd2356ee10df14ba3653d9f854f9f4b93778d0d0b6efd547d9070a962e12577055144f0696c4f56c9f2a44c7cbc83fa07e089587b\u0026X-Goog-SignedHeaders=host", nil
	} else if name == "videos/GSVsfCCtRr0/[기생충] 30초 예고_360p.mp4" {
		return "https://storage.googleapis.com/youtube-download-backend-beta/videos/GSVsfCCtRr0/test.mp4?X-Goog-Algorithm=GOOG4-RSA-SHA256\u0026X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210801%2Fauto%2Fstorage%2Fgoog4_request\u0026X-Goog-Date=20210801T084013Z\u0026X-Goog-Expires=899\u0026X-Goog-Signature=31390355ce8353e169279bc4286078c29ef2ed5bcb2bcc13380c67efac26e4155275e0d1a3683cc77e7d9188e33ffaee0e7056c584951e339c327934e7d08a78a0e549acd05b1506ed62c6e048c8500416cc7086f07bd9a03fa1df6b678220b0a1c810d0715f904c889a80fc3208a997ebb6e117f6857cf526daa32bb515c048384f85b27456c6c39ae7647efe2f0be24f0905a52fe0e5b94bbd6871a71c34951f1317cbccee23fa4ed0566415e9b55221fba05558ebdbc9c6e48f11396e1766f7f62b49756b8daa7cc54c6169176f81e6726e91e02dec30d671d7cb4cca85cc6556c5ef312293ece4b22792b53aa05481e3ad0ba225ca3435e935adb25a7143\u0026X-Goog-SignedHeaders=host", nil
	} else {
		return "", errors.New(fmt.Sprintf("could not mock %s", name))
	}
}
