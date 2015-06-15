all:
	go build
assets:
	go-bindata templates public public/css public/css/components public/fonts public/js public/js/components public/js/core
