all: linux


validate:
	@go build ./
	@go vet ./
	@go tool vet -shadow ./
	@golint ./
	@ineffassign ./

linux:
	@env GOOS=linux GOARCH=amd64 go install -ldflags "-X main.buildtime=`date -u +%Y-%m-%dT%H:%M:%S%z` -X main.githash=`git rev-parse HEAD` -X main.shorthash=`git rev-parse --short HEAD` -X main.builddate=`date -u +%Y%m%d`" ./
