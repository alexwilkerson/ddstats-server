.PHONY: clean build

clean:
	rm -f ./dist/server
	rm -f ./dist/collector

build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/server -v ./cmd/server
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/collector -v ./cmd/collector

deploy: build
	scp ./dist/server casd:~/ddstats
	scp ./dist/collector casd:~/ddstats
