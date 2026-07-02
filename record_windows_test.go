//go:build windows && amd64

package mvs

import (
	"errors"
	"testing"
)

func TestNormalizeRecordOptionsValid(t *testing.T) {
	c := &Camera{}
	options, err := c.normalizeRecordOptions(RecordOptions{
		Path:        "record.avi",
		PixelType:   PixelTypeMono8,
		Width:       96,
		Height:      96,
		FrameRate:   25,
		BitRateKbps: 128,
		Format:      RecordFormatAVI,
	})
	if err != nil {
		t.Fatal(err)
	}
	if options.Format != RecordFormatAVI || options.FrameRate != 25 || options.BitRateKbps != 128 {
		t.Fatalf("unexpected normalized options: %+v", options)
	}
}

func TestNormalizeRecordOptionsRejectsInvalidSize(t *testing.T) {
	c := &Camera{}
	_, err := c.normalizeRecordOptions(RecordOptions{
		Path:        "record.avi",
		PixelType:   PixelTypeMono8,
		Width:       97,
		Height:      96,
		FrameRate:   25,
		BitRateKbps: 128,
		Format:      RecordFormatAVI,
	})
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("error=%v, want ErrInvalidArgument", err)
	}
}

func TestNormalizeRecordOptionsRejectsInvalidRate(t *testing.T) {
	c := &Camera{}
	_, err := c.normalizeRecordOptions(RecordOptions{
		Path:        "record.avi",
		PixelType:   PixelTypeMono8,
		Width:       96,
		Height:      96,
		FrameRate:   2000,
		BitRateKbps: 128,
		Format:      RecordFormatAVI,
	})
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("error=%v, want ErrInvalidArgument", err)
	}
}
