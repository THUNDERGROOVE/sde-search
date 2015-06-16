all:
	go build -v -ldflags "-X main.Version `git rev-parse --short HEAD` -X main.Branch `git rev-parse --abbrev-ref HEAD`"
install:
	go install -v -ldflags "-X main.Version `git rev-parse --short HEAD` -X main.Branch `git rev-parse --abbrev-ref HEAD`"
assets:
	go-bindata templates public public/css public/css/components public/fonts public/js public/js/components public/js/core