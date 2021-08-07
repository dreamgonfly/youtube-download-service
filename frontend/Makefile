install:
	npm install

start:
	REACT_APP_ENV=local npm start

build-beta:
	REACT_APP_ENV=beta npm run build

build-live:
	REACT_APP_ENV=live npm run build

deploy-beta:
	gsutil rsync -R build/ gs://youtube-download-frontend-beta

deploy-live:
	gsutil rsync -R build/ gs://youtube-download-frontend-live
