.PHONY: build doc
build:
	go build ./...

doc:
	go get github.com/GandalfUK/godoc2ghmd
	go generate ./...
