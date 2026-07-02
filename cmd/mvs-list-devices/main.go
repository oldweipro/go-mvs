//go:build windows && amd64

package main

import (
	"fmt"
	"log"

	"github.com/oldweipro/go-mvs"
)

func main() {
	sdk, err := mvs.New(mvs.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := sdk.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer sdk.Finalize()

	fmt.Printf("MVS SDK: %s\n", sdk.Version())
	fmt.Printf("DLL: %s\n", sdk.DLLPath())

	devices, err := sdk.EnumerateDefaultDevices()
	if err != nil {
		log.Fatal(err)
	}
	if len(devices) == 0 {
		fmt.Println("no devices found")
		return
	}

	for _, device := range devices {
		fmt.Printf(
			"[%d] %s model=%q serial=%q user=%q ip=%q\n",
			device.Index,
			device.TransportLayerName,
			device.ModelName,
			device.SerialNumber,
			device.UserDefinedName,
			device.CurrentIP,
		)
	}
}
