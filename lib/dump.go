package lib

import (
	"encoding/csv"
	gt7 "github.com/snipem/go-gt7-telemetry/lib"
	"log"
	"os"
	"strconv"
	"time"
)

func readCsv(filename string) (gtdata []gt7.GTData) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Unable to read input file "+filename, err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {

		packageId, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}
		rpm, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			log.Fatal(err)
		}
		throttle, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			log.Fatal(err)
		}
		brake, err := strconv.ParseFloat(record[3], 32)
		if err != nil {
			log.Fatal(err)
		}
		gear, err := strconv.Atoi(record[4])
		if err != nil {
			log.Fatal(err)
		}

		gtdata = append(gtdata, gt7.GTData{
			PackageID:   int32(packageId),
			RPM:         float32(rpm),
			Throttle:    float32(throttle),
			Brake:       float32(brake),
			CurrentGear: uint8(gear),
		})

	}

	return gtdata
}

type GT7Dump struct {
	LastData gt7.GTData
	data     []gt7.GTData
}

func NewGT7Dump(filename string) GT7Dump {
	data := readCsv(filename)
	gt7d := GT7Dump{
		LastData: gt7.GTData{},
		data:     data,
	}
	return gt7d
}

func (gt7d *GT7Dump) Run() {

	for {
		for i := 0; i < len(gt7d.data); i++ {
			gt7d.LastData = gt7d.data[i]
			time.Sleep(16 * time.Millisecond)
		}
		// Start over
	}

}
