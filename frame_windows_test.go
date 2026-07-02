//go:build windows && amd64

package mvs

import (
	"encoding/binary"
	"testing"
	"unsafe"
)

func TestFrameFromInfoCopiesSubImages(t *testing.T) {
	partA := []byte{1, 2}
	partB := []byte{3, 4}
	rawImages := []mvCCImage{
		{
			Width:     2,
			Height:    1,
			PixelType: PixelTypeMono8,
			ImageBuf:  &partA[0],
			ImageLen:  uint64(len(partA)),
		},
		{
			Width:     2,
			Height:    1,
			PixelType: PixelTypeMono8,
			ImageBuf:  &partB[0],
			ImageLen:  uint64(len(partB)),
		},
	}
	info := mvFrameOutInfoEx{
		Width:        2,
		Height:       2,
		PixelType:    PixelTypeMono8,
		FrameLen:     4,
		ExtraType:    FrameExtraSubImages,
		SubImageNum:  uint32(len(rawImages)),
		SubImageList: unsafe.Pointer(&rawImages[0]),
	}

	frame, err := frameFromInfo(&info, []byte{1, 2, 3, 4})
	if err != nil {
		t.Fatal(err)
	}
	if frame.ExtraType != FrameExtraSubImages || frame.SubImageNum != 2 {
		t.Fatalf("unexpected extra info: type=%d num=%d", frame.ExtraType, frame.SubImageNum)
	}
	if len(frame.Parts) != 2 {
		t.Fatalf("parts=%d, want 2", len(frame.Parts))
	}
	if got := frame.Parts[0].Data; len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("part[0] data=%v", got)
	}
	partA[0] = 9
	if frame.Parts[0].Data[0] != 1 {
		t.Fatal("subimage data was not copied")
	}
}

func TestFrameFromInfoCopiesMultiParts(t *testing.T) {
	partData := []byte{9, 8, 7, 6}
	var specific mvGigePartDataInfo
	binary.LittleEndian.PutUint32(specific.Data[0:4], 640)
	binary.LittleEndian.PutUint32(specific.Data[4:8], 480)

	rawParts := []mvGigeMultiPartInfo{
		{
			DataType:         uint32(MultiPartData3DImage1Planar),
			DataFormat:       PixelTypeCoord3DABC32f,
			SourceID:         1,
			RegionID:         2,
			DataPurposeID:    3,
			Zones:            1,
			Length:           uint64(len(partData)),
			PartAddr:         &partData[0],
			DataTypeSpecific: specific,
		},
	}
	info := mvFrameOutInfoEx{
		Width:        640,
		Height:       480,
		PixelType:    PixelTypeCoord3DABC32f,
		FrameLen:     uint32(len(partData)),
		ExtraType:    FrameExtraMultiParts,
		SubImageNum:  uint32(len(rawParts)),
		SubImageList: unsafe.Pointer(&rawParts[0]),
	}

	frame, err := frameFromInfo(&info, partData)
	if err != nil {
		t.Fatal(err)
	}
	if len(frame.Parts) != 1 {
		t.Fatalf("parts=%d, want 1", len(frame.Parts))
	}
	part := frame.Parts[0]
	if part.DataType != MultiPartData3DImage1Planar || part.PixelType != PixelTypeCoord3DABC32f {
		t.Fatalf("unexpected multipart metadata: %+v", part)
	}
	if part.Width != 640 || part.Height != 480 {
		t.Fatalf("multipart size=%dx%d, want 640x480", part.Width, part.Height)
	}
	if got := part.Data; len(got) != 4 || got[0] != 9 || got[3] != 6 {
		t.Fatalf("multipart data=%v", got)
	}
	partData[0] = 0
	if frame.Parts[0].Data[0] != 9 {
		t.Fatal("multipart data was not copied")
	}
}
