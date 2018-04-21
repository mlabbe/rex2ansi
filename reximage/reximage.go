package reximage

import (
	"os"
	"log"
	"fmt"
	"errors"
	"compress/gzip"
	"encoding/binary"
)

type RexRGB struct {
	Red   byte
	Green byte
	Blue  byte
}

type RexCell struct {
	Glyph byte
	_     [3]byte
	Fg    RexRGB
	Bg    RexRGB
}

type RexLayer struct {
	Width   int32
	Height  int32
	Cells   []RexCell
}

type RexImage struct {
	Version    uint32
	LayerCount uint32
	Layers     []RexLayer
}

// in-place flattening of rexpaint layers
func (srcImg *RexImage) Flatten() {
	cellCount :=
		srcImg.Layers[0].Width *
		srcImg.Layers[0].Height

	dstLayer := RexLayer{
		Width: srcImg.Layers[0].Width,
		Height: srcImg.Layers[0].Height,
		Cells: make([]RexCell, cellCount)}

	// initialize the dstLayer to all transparent to start
	for j := 0; j < int(cellCount); j++ {
		dstLayer.Cells[j] = RexCell{
			Glyph: 0x0,
			Fg:RexRGB{Red:0,   Green:0, Blue:0},
			Bg:RexRGB{Red:255, Green:0, Blue:255}}
	}

	// Apply lowset-to-highest
	for i := 0; i < int(srcImg.LayerCount); i++ {
		srcLayer := &srcImg.Layers[i]

		for j := 0; j < int(cellCount); j++ {
			// assign non-transparent cells only
			if !srcLayer.Cells[j].IsTransparent() {
				dstLayer.Cells[j] = srcLayer.Cells[j]
			}
		}
	}

	srcImg.Layers = []RexLayer{dstLayer}
	srcImg.LayerCount = 1
}

// check for colorkey transparency (value from rexpaint)
func (c RexCell) IsTransparent() bool {
	return c.Bg.Red == 255 &&
		c.Bg.Green == 0 &&
		c.Bg.Blue == 255
}

func CompareRGB(a, b RexRGB) bool {
	return a.Red == b.Red &&
		a.Green == b.Green &&
		a.Blue  == b.Blue
}

// Read a rexpaint file into a RexImage structure
func Read(file *os.File, verbose bool) (*RexImage, error) {

	// open gzip
	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	var image RexImage
	err = binary.Read(reader, binary.LittleEndian, &image.Version)
	if err != nil {
		return nil, errors.New("Could not read verison")
	}

	err = binary.Read(reader, binary.LittleEndian, &image.LayerCount)
	if err != nil {
		return nil, errors.New("Could not read layer count")
	}

	if verbose {
		log.Printf("Version: %d 0x%x", image.Version, image.Version)
		log.Printf("NumLayers: %d", image.LayerCount)
	}

	image.Layers = make([]RexLayer, image.LayerCount)

	for i := 0; i < int(image.LayerCount); i++ {
		layer := &image.Layers[i]

		err = binary.Read(reader, binary.LittleEndian, &layer.Width)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Could not read width on layer %d", i))
		}

		err = binary.Read(reader, binary.LittleEndian, &layer.Height)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Could not read height on layer %d", i))
		}

		cellCount := layer.Width * layer.Height
		if cellCount <= 0 {
			return nil, errors.New(
				fmt.Sprintf("Invalid cell count on layer %d: %d (%xx%x)",
					i, cellCount, layer.Width, layer.Height))
		}
		layer.Cells = make([]RexCell, cellCount)

		for j := 0; j < int(cellCount); j++ {
			c := RexCell{}
			err := binary.Read(reader, binary.LittleEndian, &c)
			if err != nil {
				return nil, errors.New("Could not read cell data")
			}

			// transpose column-major to row-major
			k := int32(j)
			x, y := k/layer.Height, k%layer.Height
			layer.Cells[y*layer.Width + x] = c
		}

		if verbose {
			log.Printf("layer %d dims: %dx%d\n", i, layer.Width, layer.Height)
		}
	}

	return &image, nil
}
