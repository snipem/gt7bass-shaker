package lib

import (
	"fmt"
	"testing"
)

func TestGT7Dump_Run(t *testing.T) {
	gt7c := NewGT7Dump("../dump.csv")
	go gt7c.Run()

	lastPackageId := int32(0)
	for {
		if lastPackageId != gt7c.LastData.PackageID {
			fmt.Println(gt7c.LastData.RPM)
		}
	}
}
