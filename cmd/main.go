// demo package simulating a realtime generation and processing.
// Start the example from your terminal and type a letter + enter.
package main

import (
	"flag"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/generator"
	"github.com/go-audio/transforms"
	"github.com/gordonklaus/portaudio"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"github.com/snipem/gt7buttkicker/lib"
	"log"
	"os"
)

func main() {

	var inputDumpFile string
	flag.StringVar(&inputDumpFile, "input-dump-file", "", "Specifies the input dump file, will use PlayStation if not set")

	flag.Parse() // parse the flags from the command line, see https://golang.org/pkg/flag/i
	fmt.Println(os.Args)

	if inputDumpFile == "" {
		fmt.Println("Using PlayStation as telemetry input")
		gt7c := gt7.NewGT7Communication("255.255.255.255")
		go gt7c.Run()
		Play(&gt7c.LastData)
	} else {
		fmt.Println("Using dump file as telemetry input")
		gt7c := lib.NewGT7Dump(inputDumpFile)
		go gt7c.Run()
		Play(&gt7c.LastData)
	}

	fmt.Println("done")
}

func Play(ld *gt7.GTData) {
	currentNote := 440.0

	rpmGenerator := generator.NewOsc(generator.WaveSine, currentNote, audio.FormatMono44100.SampleRate)
	rpmGenerator.Amplitude = 1

	brakeGenerator := generator.NewOsc(generator.WaveSine, currentNote, audio.FormatMono44100.SampleRate)
	brakeGenerator.Amplitude = 1

	gainControl := 0.0
	currentVol := float64(1)

	bufferSize := 512

	rpmbuf := &audio.FloatBuffer{
		Data:   make([]float64, bufferSize),
		Format: audio.FormatMono44100,
	}

	brakebuf := &audio.FloatBuffer{
		Data:   make([]float64, bufferSize),
		Format: audio.FormatMono44100,
	}

	buf := &audio.FloatBuffer{
		Data:   make([]float64, bufferSize),
		Format: audio.FormatMono44100,
	}

	go func() {
		// track gt7
		oldPackageId := int32(-1)

		for {
			if ld.PackageID != oldPackageId {
				fmt.Println(oldPackageId)

				//if gt7c.LastData.Brake > 0 {
				//	fmt.Println("Brake")
				brakeGenerator.SetFreq(float64(ld.Brake + 32))
				//} else {
				//	fmt.Println("RPM")
				rpmGenerator.SetFreq(float64(ld.RPM) / 28)
				//}
			}
			oldPackageId = ld.PackageID
		}
	}()

	//buf2 := &audio.FloatBuffer{
	//	Data:   make([]float64, bufferSize),
	//	Format: audio.FormatMono44100,
	//}
	// Audio output
	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]float32, bufferSize)
	stream, err := portaudio.OpenDefaultStream(0, 2, 44100, len(out), &out)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()
	for {

		//transform.NormalizeMax()

		// populate the out buffer
		if err := rpmGenerator.Fill(rpmbuf); err != nil {
			log.Printf("error filling up the buffer")
		}
		if err := brakeGenerator.Fill(brakebuf); err != nil {
			log.Printf("error filling up the buffer")
		}
		// apply vol control if needed (applied as a transform instead of a control
		// on the osc)
		if gainControl != 0 {
			currentVol += gainControl
			if currentVol < 0.1 {
				currentVol = 0
			}
			if currentVol > 6 {
				currentVol = 6
			}
			fmt.Printf("new vol %f.2", currentVol)
			gainControl = 0
		}

		//bufout := mix(buf, buf2)

		// chose buffer
		if ld.Brake > 0 {
			buf = brakebuf
		} else {
			buf = rpmbuf
		}

		transforms.Gain(buf, currentVol)

		f64ToF32Copy(out, buf.Data)

		// write to the stream
		if err := stream.Write(); err != nil {
			log.Printf("error writing to stream : %v\n", err)
		}
	}
}

func mix(buf *audio.FloatBuffer, buf2 *audio.FloatBuffer) (mixBuf *audio.FloatBuffer) {
	bufferSize := 512
	mixBuf = &audio.FloatBuffer{
		Data:   make([]float64, bufferSize),
		Format: audio.FormatMono44100,
	}

	for i, _ := range mixBuf.Data {
		mixBuf.Data[i] = buf.Data[i] + buf2.Data[i]
	}

	transforms.NormalizeMax(mixBuf)
	return mixBuf
}

func switchNote(data float32, osc *generator.Osc) {
	currentNote := float64(data / 28)
	osc.SetFreq(currentNote)
}

// portaudio doesn't support float64 so we need to copy our data over to the
// destination buffer.
func f64ToF32Copy(dst []float32, src []float64) {
	for i := range src {
		dst[i] = float32(src[i])
	}
}
