build:

deps:
	brew install arm-linux-gnueabihf-binutils

build_rasp_on_mac:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -o gt7buttkicker.arm.bin cmd/main.go

build_raspberrypi:
	apt-get update
	apt-get install pkg-config
	go mod tidy
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -o gt7buttkicker.arm.bin cmd/main.go

build_using_docker:
	docker run -it --rm \
	  -v /home/matze/work/gt7buttkicker:/go/src/github.com/user/go-project \
	  -w /go/src/github.com/user/go-project \
	  -e CGO_ENABLED=1 \
	  docker.elastic.co/beats-dev/golang-crossbuild:1.20-armhf \
	  --build-cmd "make build_raspberrypi" \
	  -p "linux/armv7"
