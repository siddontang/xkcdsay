all: build

build:
	go build -o bin/xkcdsay ./cmd/xkcdsay/.
	go build -o bin/xkcddown ./cmd/xkcddown/.
	go build -o bin/xkcdsync ./cmd/xkcdsync/.
	GOOS=linux go build -o bin/xkcd_lambda_sync ./cmd/xkcd_lambda_sync/.