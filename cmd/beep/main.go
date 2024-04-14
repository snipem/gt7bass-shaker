package main

import (
	"embed"
	"fmt"
	"github.com/gopxl/beep/wav"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"github.com/snipem/gt7tools/lib/dump"
	"log"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
)

//go:embed wav/knock_short.wav
var embedFile embed.FS

func shift() {
	fwav, err := embedFile.Open("wav/knock_short.wav")
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := wav.Decode(fwav)
	if err != nil {
		log.Fatal(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	ctrl := &beep.Ctrl{Streamer: beep.Loop(1, streamer), Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	speedy := beep.ResampleRatio(4, 1, volume)
	speaker.Play(speedy)
}

func main() {

	gt7c := gt7.NewGT7Communication("255.255.255.255")

	dumpFilePath := ""
	if len(os.Args) > 1 {
		dumpFilePath = os.Args[1]
	}

	if dumpFilePath != "" {

		gt7dump, err := dump.NewGT7Dump(dumpFilePath, gt7c)
		if err != nil {
			log.Fatalf("Error loading dump file: %v", err)
		}
		log.Println("Using dump file: ", dumpFilePath)
		go gt7dump.Run()

	} else {
		go func() {

			for {
				err := gt7c.Run()
				if err != nil {
					log.Printf("error running gt7c.Run(): %v", err)
				}
				log.Println("Sleeping 10 seconds before restarting gt7c.Run()")
				time.Sleep(10 * time.Second)
			}
		}()
	}
	//f, err := os.Open("wav/Miami_Slice_-_04_-_Step_Into_Me.mp3")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//streamer, format, err := mp3.Decode(f)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer streamer.Close()
	//
	//speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	//
	//ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	//volume := &effects.Volume{
	//	Streamer: ctrl,
	//	Base:     2,
	//	Volume:   0,
	//	Silent:   false,
	//}
	//beep.ResampleRatio(4, 1, volume)
	//speaker.Play(speedy)

	oldPackageId := int32(0)
	oldGear := uint8(0)

	for {
		//fmt.Print("Press [ENTER] to pause/resume. ")
		//fmt.Scanln()

		if oldPackageId != gt7c.LastData.PackageID {
			//fmt.Printf("%d: %d\n", gt7c.LastData.PackageID, gt7c.LastData.CurrentGear)
			if oldGear != gt7c.LastData.CurrentGear {
				speaker.Lock()
				//ctrl.Paused = !ctrl.Paused
				fmt.Printf("%d: knock %d -> %d\n", gt7c.LastData.PackageID, oldGear, gt7c.LastData.CurrentGear)
				go shift()
				//volume.Volume += 0.5
				//speedy.SetRatio(speedy.Ratio() + 0.1)
				speaker.Unlock()
			}
			oldGear = gt7c.LastData.CurrentGear
		}
		oldPackageId = gt7c.LastData.PackageID
	}
}
