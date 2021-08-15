install:
	npm install

start:
	REACT_APP_ENV=local npm start

build-beta:
	REACT_APP_ENV=beta npm run build

build-live:
	REACT_APP_ENV=live npm run build

deploy-beta:
	# Exclude .DS_Store: https://stackoverflow.com/a/35210687
	gsutil rsync -rx \..* build/ gs://youtube-download-frontend-beta

deploy-live:
	gsutil rsync -rx \..* build/ gs://youtube-download-frontend-live
