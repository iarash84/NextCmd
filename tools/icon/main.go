// Command icon prepares a transparent PNG and converts it into a
// multi-resolution Windows ICO file.
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

var sizes = []int{16, 20, 24, 32, 40, 48, 64, 128, 256}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: go run ./tools/icon <input.png> <output.png> <output.ico>")
		os.Exit(2)
	}
	input, err := os.Open(os.Args[1])
	if err != nil {
		fatal(err)
	}
	defer input.Close()
	decoded, err := png.Decode(input)
	if err != nil {
		fatal(err)
	}
	source := removeLightBackground(decoded)

	pngOutput, err := os.Create(os.Args[2])
	if err != nil {
		fatal(err)
	}
	if err := png.Encode(pngOutput, source); err != nil {
		_ = pngOutput.Close()
		fatal(err)
	}
	if err := pngOutput.Close(); err != nil {
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

	output, err := os.Create(os.Args[3])
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
			x0 := bounds.Min.X + x*bounds.Dx()/size
			x1 := bounds.Min.X + (x+1)*bounds.Dx()/size
			y0 := bounds.Min.Y + y*bounds.Dy()/size
			y1 := bounds.Min.Y + (y+1)*bounds.Dy()/size
			if x1 <= x0 {
				x1 = x0 + 1
			}
			if y1 <= y0 {
				y1 = y0 + 1
			}
			var red, green, blue, alpha, count uint64
			for sy := y0; sy < y1; sy++ {
				for sx := x0; sx < x1; sx++ {
					pixel := color.NRGBAModel.Convert(source.At(sx, sy)).(color.NRGBA)
					red += uint64(pixel.R) * uint64(pixel.A)
					green += uint64(pixel.G) * uint64(pixel.A)
					blue += uint64(pixel.B) * uint64(pixel.A)
					alpha += uint64(pixel.A)
					count++
				}
			}
			pixel := color.NRGBA{A: uint8(alpha / count)}
			if alpha > 0 {
				pixel.R = uint8(red / alpha)
				pixel.G = uint8(green / alpha)
				pixel.B = uint8(blue / alpha)
			}
			destination.SetNRGBA(x, y, pixel)
		}
	}
	return destination
}

// removeLightBackground converts the light checkerboard emitted by some image
// editors into real alpha. The NextCmd artwork itself contains no near-white
// pixels, so this preserves the mark and its dark tile while smoothing the edge.
func removeLightBackground(source image.Image) *image.NRGBA {
	bounds := source.Bounds()
	result := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			pixel := color.NRGBAModel.Convert(source.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.NRGBA)
			minimum := min(pixel.R, pixel.G, pixel.B)
			switch {
			case minimum >= 220:
				pixel.A = 0
			case minimum > 180:
				pixel.A = uint8((int(220-minimum) * int(pixel.A)) / 40)
			}
			result.SetNRGBA(x, y, pixel)
		}
	}
	return result
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "icon:", err)
	os.Exit(1)
}
