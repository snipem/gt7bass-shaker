build_on_raspberrypi:
	git pull
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o gt7buttkicker.arm.bin cmd/main.go

build_rasp_on_mac:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -o gt7buttkicker.arm.bin cmd/main.go

deps:
	brew install arm-linux-gnueabihf-binutils

deps_on_rasbperrypi:
	apt-get update
	apt-get install git pkg-config portaudio19-dev

build_using_docker:
	docker run -it --rm \
	  -v /home/matze/work/gt7buttkicker:/go/src/github.com/user/go-project \
	  -w /go/src/github.com/user/go-project \
	  -e CGO_ENABLED=1 \
	  docker.elastic.co/beats-dev/golang-crossbuild:1.20-armhf \
	  --build-cmd "make build_raspberrypi" \
	  -p "linux/armv7"
