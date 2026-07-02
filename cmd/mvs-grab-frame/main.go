//go:build windows && amd64

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/oldweipro/go-mvs"
)

func main() {
	output := flag.String("out", "frame.raw", "output raw frame path")
	imageOutput := flag.String("image", "", "optional image output path: bmp, jpg, jpeg, png, tif, or tiff")
	imageFormat := flag.String("format", "auto", "image format: auto, bmp, jpeg, png, tif")
	jpegOutput := flag.String("jpeg", "", "optional JPEG image path")
	quality := flag.Uint("quality", 90, "JPEG quality, valid range is 51-99")
	serial := flag.String("serial", "", "open device by serial number")
	ip := flag.String("ip", "", "open device by current IP address")
	userName := flag.String("user", "", "open device by user-defined name")
	timeout := flag.Duration("timeout", time.Second, "frame timeout")
	flag.Parse()

	sdk, err := mvs.New(mvs.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := sdk.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer sdk.Finalize()

	camera, err := openCamera(sdk, *serial, *ip, *userName)
	if err != nil {
		log.Fatal(err)
	}
	defer camera.Close()

	if err := camera.ConfigureOptimalPacketSize(); err != nil {
		log.Printf("packet size setup skipped: %v", err)
	}
	if err := camera.SetTriggerMode(false); err != nil {
		log.Fatal(err)
	}
	if err := camera.StartGrabbing(); err != nil {
		log.Fatal(err)
	}
	defer camera.StopGrabbing()

	frame, err := camera.GetFrame(*timeout)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(*output, frame.Data, 0644); err != nil {
		log.Fatal(err)
	}
	if *imageOutput != "" && *jpegOutput != "" {
		log.Fatal("-image and -jpeg cannot be used together")
	}
	savePath := *imageOutput
	saveType := mvs.ImageTypeUndefined
	if *jpegOutput != "" {
		savePath = *jpegOutput
		saveType = mvs.ImageTypeJPEG
	}
	if savePath != "" {
		if saveType == mvs.ImageTypeUndefined {
			var err error
			saveType, err = parseImageType(*imageFormat, savePath)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := camera.SaveFrameToFile(frame, savePath, mvs.ImageSaveOptions{
			Type:    saveType,
			Quality: uint32(*quality),
		}); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf(
		"saved %s width=%d height=%d pixel=%s frame=%d bytes=%d\n",
		*output,
		frame.Width,
		frame.Height,
		frame.PixelTypeName(),
		frame.FrameNum,
		len(frame.Data),
	)
}

func parseImageType(value string, path string) (mvs.ImageType, error) {
	switch strings.ToLower(value) {
	case "", "auto":
		imageType, ok := mvs.ImageTypeFromExtension(path)
		if !ok {
			return mvs.ImageTypeUndefined, fmt.Errorf("cannot infer image format from %q", path)
		}
		return imageType, nil
	case "bmp":
		return mvs.ImageTypeBMP, nil
	case "jpg", "jpeg":
		return mvs.ImageTypeJPEG, nil
	case "png":
		return mvs.ImageTypePNG, nil
	case "tif", "tiff":
		return mvs.ImageTypeTIFF, nil
	default:
		return mvs.ImageTypeUndefined, fmt.Errorf("unsupported image format %q", value)
	}
}

func openCamera(sdk *mvs.SDK, serial string, ip string, userName string) (*mvs.Camera, error) {
	switch {
	case serial != "":
		return sdk.OpenDeviceBySerial(serial, mvs.AccessExclusive)
	case ip != "":
		return sdk.OpenDeviceByIP(ip, mvs.AccessExclusive)
	case userName != "":
		return sdk.OpenDeviceByUserDefinedName(userName, mvs.AccessExclusive)
	default:
		devices, err := sdk.EnumerateDefaultDevices()
		if err != nil {
			return nil, err
		}
		if len(devices) == 0 {
			return nil, fmt.Errorf("no devices found")
		}
		return sdk.OpenDevice(devices[0], mvs.AccessExclusive)
	}
}
