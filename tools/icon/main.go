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
	source := frameMark(removeBackground(decoded))

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

// removeBackground preserves both the cyan and deep-blue portions of the mark.
// Blue-minus-red separates the artwork from the navy tile without clipping its
// low-green gradient. Detached tile highlights are removed by component area.
// This is an artwork-specific key, not a general-purpose background remover.
func removeBackground(source image.Image) *image.NRGBA {
	bounds := source.Bounds()
	result := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			pixel := color.NRGBAModel.Convert(source.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.NRGBA)
			chroma := int(pixel.B) - int(pixel.R)
			switch {
			case chroma <= 70:
				pixel = color.NRGBA{}
			case chroma < 150:
				pixel.A = uint8(int(pixel.A) * (chroma - 70) / 80)
			}
			if pixel.A == 0 {
				pixel = color.NRGBA{}
			}
			result.SetNRGBA(x, y, pixel)
		}
	}
	return removeSpecks(result)
}

// removeSpecks retains substantial connected shapes, not isolated highlights
// from the old tile. Eight-way connectivity preserves diagonal edge pixels.
func removeSpecks(source *image.NRGBA) *image.NRGBA {
	bounds := source.Bounds()
	width := bounds.Dx()
	seen := make([]bool, width*bounds.Dy())
	var components [][]image.Point
	largest := 0
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < width; x++ {
			index := y*width + x
			if seen[index] || source.NRGBAAt(x, y).A == 0 {
				continue
			}
			seen[index] = true
			points := []image.Point{{X: x, Y: y}}
			for head := 0; head < len(points); head++ {
				p := points[head]
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						n := p.Add(image.Pt(dx, dy))
						if !n.In(bounds) {
							continue
						}
						i := n.Y*width + n.X
						if !seen[i] && source.NRGBAAt(n.X, n.Y).A > 0 {
							seen[i] = true
							points = append(points, n)
						}
					}
				}
			}
			largest = max(largest, len(points))
			components = append(components, points)
		}
	}
	for _, points := range components {
		if len(points) < max(1, largest/5) {
			for _, p := range points {
				source.SetNRGBA(p.X, p.Y, color.NRGBA{})
			}
		}
	}
	return source
}

// frameMark centers the visible mark on a square canvas with an 8% safe area.
// Cropping away the old tile's margins makes the silhouette legible at 16px
// without stretching the artwork or letting it touch the icon's edges.
func frameMark(source *image.NRGBA) *image.NRGBA {
	bounds := source.Bounds()
	visible := image.Rectangle{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if source.NRGBAAt(x, y).A > 0 {
				visible = visible.Union(image.Rect(x, y, x+1, y+1))
			}
		}
	}
	if visible.Empty() {
		return source
	}
	side := max(visible.Dx(), visible.Dy())
	padding := max(1, (side*8+83)/84)
	canvas := image.NewNRGBA(image.Rect(0, 0, side+2*padding, side+2*padding))
	xOffset := (canvas.Bounds().Dx() - visible.Dx()) / 2
	yOffset := (canvas.Bounds().Dy() - visible.Dy()) / 2
	for y := 0; y < visible.Dy(); y++ {
		for x := 0; x < visible.Dx(); x++ {
			canvas.SetNRGBA(x+xOffset, y+yOffset, source.NRGBAAt(x+visible.Min.X, y+visible.Min.Y))
		}
	}
	return canvas
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "icon:", err)
	os.Exit(1)
}
