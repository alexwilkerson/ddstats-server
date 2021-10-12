.PHONY: clean build

clean:
	rm -f ./dist/server

build: clean
	GOOS=linux GOARCH=amd64 go build -o dist/server -v ./cmd/server

deploy: build
	scp ./dist/server casd:~/ddstats
