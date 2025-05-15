build:
	go build -o bin/jaabz ./cmd/.

run: build
	bin/jaabz