build:

deps:
	brew install arm-linux-gnueabihf-binutils

build_raspberrypi:
	CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -o gt7buttkicker.arm.bin cmd/main.go
