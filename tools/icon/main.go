// Command icon converts a square PNG into a multi-resolution Windows ICO file.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var sizes = []int{16, 24, 32, 48, 64, 128, 256}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run ./tools/icon <input.png> <output.ico>")
		os.Exit(2)
	}
	input, err := os.Open(os.Args[1])
	if err != nil {
		fatal(err)
	}
	defer input.Close()
	source, err := png.Decode(input)
	if err != nil {
		fatal(err)
	}

	images := make([][]byte, len(sizes))
	for i, size := range sizes {
		var encoded bytes.Buffer
		if err := png.Encode(&encoded, resize(source, size)); err != nil {
			fatal(err)
		}
		images[i] = encoded.Bytes()
	}

	output, err := os.Create(os.Args[2])
	if err != nil {
		fatal(err)
	}
	defer output.Close()
	write := func(value any) {
		if err := binary.Write(output, binary.LittleEndian, value); err != nil {
			fatal(err)
		}
	}
	write(uint16(0))
	write(uint16(1))
	write(uint16(len(images)))
	offset := uint32(6 + len(images)*16)
	for i, data := range images {
		sizeByte := uint8(sizes[i])
		if sizes[i] == 256 {
			sizeByte = 0
		}
		write(sizeByte)
		write(sizeByte)
		write(uint8(0))
		write(uint8(0))
		write(uint16(1))
		write(uint16(32))
		write(uint32(len(data)))
		write(offset)
		offset += uint32(len(data))
	}
	for _, data := range images {
		if _, err := output.Write(data); err != nil {
			fatal(err)
		}
	}
}

func resize(source image.Image, size int) *image.NRGBA {
	bounds := source.Bounds()
	destination := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			sx := bounds.Min.X + x*bounds.Dx()/size
			sy := bounds.Min.Y + y*bounds.Dy()/size
			destination.SetNRGBA(x, y, color.NRGBAModel.Convert(source.At(sx, sy)).(color.NRGBA))
		}
	}
	return destination
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "icon:", err)
	os.Exit(1)
}
