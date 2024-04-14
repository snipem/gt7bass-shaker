build_on_raspberrypi:
	git pull
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o gt7buttkicker.arm.bin cmd/main.go
	make restart_service

build_on_raspberrypi_beep:
	git pull
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o gt7buttkicker.arm.bin cmd/beep/main.go
	make restart_service

build_rasp_on_mac:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -o gt7buttkicker.arm.bin cmd/main.go

install_service:
	sudo cp gt7buttkicker.service /etc/systemd/system/
	sudo systemctl daemon-reload
	systemctl enable gt7buttkicker.service
	systemctl start gt7buttkicker.service

restart_service:
	sudo systemctl restart gt7buttkicker.service

deps:
	brew install arm-linux-gnueabihf-binutils

wavefiles:
	wget http://cd.textfiles.com/sbsw/BEEPCHMS/KLAK.WAV -O wav/knock.wav

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
