install:
	npm install

start:
	npm start

build:
	npm run build

deploy-beta:
	gsutil rsync -R build/ gs://youtube-download-frontend-beta

deploy-live:
	gsutil rsync -R build/ gs://youtube-download-frontend-live
