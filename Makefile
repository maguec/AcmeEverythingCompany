default: deps mybuild

deps:
	go get

mybuild:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build *.go

xcompile:
	goreleaser --snapshot --skip-publish --clean
