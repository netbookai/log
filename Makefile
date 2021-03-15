.PHONY: build doc
build:
	go build ./...

doc:
	go get github.com/posener/goreadme/cmd/goreadme
	goreadme -title Log -badge-godoc -import-path github.com/go-coldbrew/log -recursive > README.md
