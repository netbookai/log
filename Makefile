.PHONY: build doc
build:
	go build ./...

doc:
	go get github.com/princjef/gomarkdoc/cmd/gomarkdoc
	gomarkdoc ./... > README.md
