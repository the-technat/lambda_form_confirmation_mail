GOOS := linux
GOARCH := amd64
CGO_ENABLED := 0

deploy:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o main . 
	zip main.zip main
	aws lambda update-function-code --function-name ${FUNCTION} --zip-file fileb://main.zip
	rm -rf main.zip

vet:
	go vet ./...

lint: staticcheck
	staticcheck ./...

staticcheck: go
	go install honnef.co/go/tools/cmd/staticcheck@latest 

go:
	which go
