GUEST_BIN = ua-proxy-guest
SERVER_BIN= ua-proxy-server

PKG_URL := main
APP_VERSION = $(shell git describe --always --tags --abbrev=0 | tr -d "[v\r\n]")
DATE_VERSION := $(shell date +%Y%m%d-%H%M)
GIT_VERSION := $(shell git rev-parse --short HEAD)
GIT_DATE_VERSION := $(GIT_VERSION)-$(DATE_VERSION)
AUTO_VERSIONING := -X $(PKG_URL).Version=$(APP_VERSION) -X $(PKG_URL).BuildDate=$(DATE_VERSION) -X $(PKG_URL).CommitSHA=$(GIT_VERSION)

all: guest server

rsrc:
	command -v rsrc || go get github.com/akavel/rsrc

guest: rsrc
	rsrc -manifest ./cmd/guest/ua-proxy-guest.exe.manifest -o ./cmd/guest/rsrc.syso
	GOOS=windows GOARCH=amd64 go build -ldflags "$(AUTO_VERSIONING) -H windowsgui" -o $(GUEST_BIN).exe ./cmd/guest/

server:
	GOOS=linux go build -ldflags "$(AUTO_VERSIONING)" -o $(SERVER_BIN) ./cmd/server/

clean:
	@rm -f $(GUEST_BIN) $(SERVER_BIN)