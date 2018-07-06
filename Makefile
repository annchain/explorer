.PHONY: build 
build:
	
	go build -v -o block-browser

test:
	
	go test ./job
