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

type Channel struct {
	Generator *generator.Osc
	Buffer    *audio.FloatBuffer
	Type      string
	mix       *Mix
}

type Mix struct {
	Channels     []*Channel
	LastData     *gt7.GTData
	buf          *audio.FloatBuffer
	TCSChannel   *Channel
	RPMChannel   *Channel
	BrakeChannel *Channel
}

func (mix *Mix) NewChannel(sine generator.WaveType, channelType string) *Channel {
	c := Channel{}
	c.Type = channelType
	c.Generator = generator.NewOsc(sine, 440, audio.FormatMono44100.SampleRate)
	c.Generator.Amplitude = 1
	c.Buffer = getBuffer(512)
	c.mix = mix
	mix.Channels = append(mix.Channels, &c)

	switch channelType {
	case "RPM":
		mix.RPMChannel = &c
	case "Brake":
		mix.BrakeChannel = &c
	case "TCS":
		mix.TCSChannel = &c
	}

	return &c
}

func (c *Channel) PopulateBuffer() {
	// populate the out buffer
	if err := c.Generator.Fill(c.Buffer); err != nil {
		log.Printf("error filling up the buffer")
	}
}

func (c *Channel) SynthesizeTelemetry() {

	switch c.Type {
	case "RPM":
		c.Generator.SetFreq(float64(c.mix.LastData.RPM) / 28)
	case "Brake":
		c.Generator.SetFreq(float64(100 - c.mix.LastData.Brake + 32))
	case "TCS":
		c.Generator.SetFreq(float64(60))
	}

}

func Play(ld *gt7.GTData) {
	mix := NewMix(ld)

	//c1 := NewChannel(generator.WaveSine)

	//currentNote := 440.0
	//
	//rpmGenerator := generator.NewOsc(generator.WaveSine, currentNote, audio.FormatMono44100.SampleRate)
	//rpmGenerator.Amplitude = 1
	//
	//brakeGenerator := generator.NewOsc(generator.WaveSine, currentNote, audio.FormatMono44100.SampleRate)
	//brakeGenerator.Amplitude = 1
	//
	//tcsGenerator := generator.NewOsc(generator.WaveTriangle, currentNote, audio.FormatMono44100.SampleRate)
	//tcsGenerator.Amplitude = 1

	mix.NewChannel(generator.WaveSine, "RPM")
	mix.NewChannel(generator.WaveSine, "Brake")
	mix.NewChannel(generator.WaveSine, "TCS")

	gainControl := 0.0
	currentVol := float64(1)

	bufferSize := 512

	//rpmbuf := getBuffer(bufferSize)
	//brakeBuf := getBuffer(bufferSize)
	//
	//brakebuf := &audio.FloatBuffer{
	//	Data:   make([]float64, bufferSize),
	//	Format: audio.FormatMono44100,
	//}
	//
	//tcsBuf := &audio.FloatBuffer{
	//	Data:   make([]float64, bufferSize),
	//	Format: audio.FormatMono44100,
	//}
	//
	//buf := &audio.FloatBuffer{
	//	Data:   make([]float64, bufferSize),
	//	Format: audio.FormatMono44100,
	//}

	go func() {
		// track gt7
		oldPackageId := int32(-1)

		for {
			if ld.PackageID != oldPackageId {
				fmt.Println(oldPackageId)

				for i := 0; i < len(mix.Channels); i++ {
					mix.Channels[i].SynthesizeTelemetry()
				}

				//if gt7c.LastData.Brake > 0 {
				//	fmt.Println("Brake")
				//brakeChannel.Generator.SetFreq(float64(100 - ld.Brake + 32))
				//} else {
				//	fmt.Println("RPM")
				//rpmChannel.Generator.SetFreq(float64(ld.RPM) / 28)

				//tcsChannel.Generator.SetFreq(float64(60))
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

		for i := 0; i < len(mix.Channels); i++ {
			mix.Channels[i].PopulateBuffer()
		}

		//if err := rpmGenerator.Fill(rpmbuf); err != nil {
		//	log.Printf("error filling up the buffer")
		//}
		//if err := brakeGenerator.Fill(brakebuf); err != nil {
		//	log.Printf("error filling up the buffer")
		//}
		//if err := tcsGenerator.Fill(tcsBuf); err != nil {
		//	log.Printf("error filling up the buffer")
		//}
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

		if ld.InRace && ld.IsPaused {
			currentVol = 0
		} else {
			currentVol = 1
		}

		// block all effects
		currentVol = 0

		buf := mix.GetMixedBuffer(ld, currentVol)

		f64ToF32Copy(out, buf.Data)

		// write to the stream
		if err := stream.Write(); err != nil {
			log.Printf("error writing to stream : %v\n", err)
		}

		//time.Sleep(16 * time.Millisecond)
	}
}

func NewMix(ld *gt7.GTData) Mix {
	mix := Mix{}
	mix.LastData = ld
	mix.buf = getBuffer(512)
	return mix
}

func (mix *Mix) GetMixedBuffer(ld *gt7.GTData, currentVol float64) *audio.FloatBuffer {
	// chose buffer
	if ld.IsTCSEngaged {
		currentVol = 1
		mix.buf = mix.TCSChannel.Buffer
	} else if ld.Brake > 0 {
		currentVol = 1
		mix.buf = mix.BrakeChannel.Buffer
	} else {
		mix.buf = mix.RPMChannel.Buffer
		currentVol = 1.5
	}

	transforms.Gain(mix.buf, currentVol)
	return mix.buf
}

func getBuffer(bufferSize int) *audio.FloatBuffer {
	rpmbuf := &audio.FloatBuffer{
		Data:   make([]float64, bufferSize),
		Format: audio.FormatMono44100,
	}
	return rpmbuf
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
