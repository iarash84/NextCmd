package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestRemoveBackground(t *testing.T) {
	pixels := []struct {
		name string
		in   color.NRGBA
		want color.NRGBA
	}{
		{"navy tile", color.NRGBA{R: 17, G: 32, B: 62, A: 255}, color.NRGBA{}},
		{"white checkerboard", color.NRGBA{R: 254, G: 254, B: 254, A: 255}, color.NRGBA{}},
		{"cyan mark", color.NRGBA{G: 180, B: 254, A: 255}, color.NRGBA{G: 180, B: 254, A: 255}},
		{"deep blue mark", color.NRGBA{G: 70, B: 230, A: 255}, color.NRGBA{G: 70, B: 230, A: 255}},
		{"soft edge", color.NRGBA{G: 60, B: 110, A: 200}, color.NRGBA{G: 60, B: 110, A: 100}},
		{"transparent", color.NRGBA{G: 180, B: 254}, color.NRGBA{}},
	}
	for _, tt := range pixels {
		t.Run(tt.name, func(t *testing.T) {
			source := image.NewNRGBA(image.Rect(7, 9, 8, 10))
			source.SetNRGBA(7, 9, tt.in)
			if got := removeBackground(source).NRGBAAt(0, 0); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveBackgroundKeepsGradientAndDiscardsSpecks(t *testing.T) {
	source := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	for y := 25; y < 75; y++ {
		for x := 25; x < 75; x++ {
			source.SetNRGBA(x, y, color.NRGBA{G: uint8(70 + y), B: 250, A: 255})
		}
	}
	source.SetNRGBA(50, 1, color.NRGBA{G: 180, B: 250, A: 255})
	clean := removeBackground(source)
	if clean.NRGBAAt(50, 1).A != 0 {
		t.Fatal("detached highlight would distort the crop")
	}
	for y := 25; y < 75; y++ {
		for x := 25; x < 75; x++ {
			if clean.NRGBAAt(x, y) != source.NRGBAAt(x, y) {
				t.Fatal("logo gradient was altered")
			}
		}
	}
	if got := frameMark(clean).Bounds().Dx(); got != 60 {
		t.Fatalf("incorrect crop after removing residue: %d", got)
	}
}

func TestFrameMark(t *testing.T) {
	source := image.NewNRGBA(image.Rect(0, 0, 20, 20))
	for y := 4; y < 16; y++ {
		for x := 7; x < 13; x++ {
			source.SetNRGBA(x, y, color.NRGBA{G: 180, B: 254, A: 255})
		}
	}
	framed := frameMark(source)
	if framed.Bounds().Dx() != 16 || framed.Bounds().Dy() != 16 {
		t.Fatalf("unexpected canvas: %v", framed.Bounds())
	}
	if framed.NRGBAAt(5, 2).A != 255 || framed.NRGBAAt(10, 13).A != 255 || framed.NRGBAAt(4, 2).A != 0 {
		t.Fatal("mark is not centered or its aspect ratio changed")
	}
	if got := frameMark(image.NewNRGBA(image.Rect(0, 0, 20, 20))); got.Bounds().Dx() != 20 {
		t.Fatal("empty artwork must remain valid")
	}
}

func TestResizePreservesEdgeColor(t *testing.T) {
	source := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	source.SetNRGBA(0, 0, color.NRGBA{G: 180, B: 254, A: 255})
	got := resize(source, 1).NRGBAAt(0, 0)
	if got != (color.NRGBA{G: 180, B: 254, A: 63}) {
		t.Fatalf("transparent pixels polluted the edge color: %v", got)
	}
}

func TestCommittedArtwork(t *testing.T) {
	data, err := os.ReadFile("../../assets/nextcmd.png")
	if err != nil {
		t.Fatal(err)
	}
	artwork, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	assertTransparentMark(t, artwork)
	ico, err := os.ReadFile("../../assets/nextcmd.ico")
	if err != nil {
		t.Fatal(err)
	}
	if len(ico) < 6+16*len(sizes) || binary.LittleEndian.Uint16(ico[2:4]) != 1 || int(binary.LittleEndian.Uint16(ico[4:6])) != len(sizes) {
		t.Fatal("invalid multi-resolution ICO header")
	}
	for i, size := range sizes {
		entry := ico[6+i*16 : 6+(i+1)*16]
		length := int(binary.LittleEndian.Uint32(entry[8:12]))
		offset := int(binary.LittleEndian.Uint32(entry[12:16]))
		if offset < 6+16*len(sizes) || length <= 0 || offset > len(ico)-length {
			t.Fatalf("invalid ICO entry for %dpx", size)
		}
		frame, err := png.Decode(bytes.NewReader(ico[offset : offset+length]))
		if err != nil {
			t.Fatal(err)
		}
		if frame.Bounds().Dx() != size || frame.Bounds().Dy() != size {
			t.Fatalf("incorrect frame dimensions for %dpx", size)
		}
		assertTransparentMark(t, frame)
	}
}

func assertTransparentMark(t *testing.T, artwork image.Image) {
	t.Helper()
	visible := 0
	bounds := artwork.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			p := color.NRGBAModel.Convert(artwork.At(x, y)).(color.NRGBA)
			if (x == bounds.Min.X || y == bounds.Min.Y || x == bounds.Max.X-1 || y == bounds.Max.Y-1) && p.A != 0 {
				t.Fatal("icon does not have transparent padding")
			}
			if p.A > 0 {
				visible++
				if int(p.B)-int(p.R) <= 70 {
					t.Fatalf("background remains at (%d,%d): %v", x, y, p)
				}
			}
		}
	}
	if visible == 0 {
		t.Fatal("icon is empty")
	}
	if visible*100 < bounds.Dx()*bounds.Dy()*20 {
		t.Fatal("mark occupies too little of the canvas; check for detached residue")
	}
}
