build:
	go build -v

tidy:
	go mod tidy

test: 
	go test -cover -parallel 5 -failfast  ./...
