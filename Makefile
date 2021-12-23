all: build

build:
	go build -o bin/xkcdsay ./cmd/xkcdsay/.
	go build -o bin/xkcddown ./cmd/xkcddown/.
	go build -o bin/xkcdsync ./cmd/xkcdsync/.