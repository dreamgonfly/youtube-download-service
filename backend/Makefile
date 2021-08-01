build-local:
	docker build --tag youtube-download-backend .

run-local:
	docker run --rm --publish 8000:8000 -e GOOGLE_APPLICATION_CREDENTIALS=/credentials/youtube-download-service-credentials.json youtube-download-backend

build-beta:
	docker build --tag asia.gcr.io/youtube-download-service/youtube-download-backend:beta .

push-beta:
	docker push asia.gcr.io/youtube-download-service/youtube-download-backend:beta

deploy-beta:
	gcloud run deploy youtube-download-backend-beta --concurrency 1 --cpu 2 --memory 8Gi --timeout 60m --image asia.gcr.io/youtube-download-service/youtube-download-backend:beta --update-env-vars ENV=beta,GOOGLE_APPLICATION_CREDENTIALS=/credentials/youtube-download-service-credentials.json --allow-unauthenticated --region asia-northeast3 --use-http2 --service-account youtube-download-service@youtube-download-service.iam.gserviceaccount.com

build-live:
	docker build --tag asia.gcr.io/youtube-download-service/youtube-download-backend:live .

push-live:
	docker push asia.gcr.io/youtube-download-service/youtube-download-backend:live

deploy-live:
	gcloud run deploy youtube-download-backend-live --concurrency 1 --cpu 2 --memory 8Gi --timeout 60m --image asia.gcr.io/youtube-download-service/youtube-download-backend:live --update-env-vars ENV=live,GOOGLE_APPLICATION_CREDENTIALS=/credentials/youtube-download-service-credentials.json --allow-unauthenticated --region asia-northeast3 --use-http2 --service-account youtube-download-service@youtube-download-service.iam.gserviceaccount.com
