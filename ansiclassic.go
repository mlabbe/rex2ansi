package main

import (
	"fmt"
	"os"

	"github.com/mlabbe/rex2ansi/reximage"
)

// ansiclassic is the catch-all phrase for codepage 437 "ASCII" ansi,
// which is what was used on DOS.  In order to view these files you
// must have codepage 437 installed on your terminal with a font that
// has the codepage 437 charset.  This is not modern, but it is what
// you will need to display these art files in a Windows console.

func classicReset() []byte {
	return []byte{'\x1b', '[', '0', 'm'}
}

func classicFg24(red, green, blue byte) []byte {
	b := []byte{'\x1b', '['}
	s := fmt.Sprintf("38;2;%d;%d;%dm", red, green, blue)
	b2 := []byte(s)
	return append(b[:], b2...)
}

func classicBg24(red, green, blue byte) []byte {
	b := []byte{'\x1b', '['}
	s := fmt.Sprintf("48;2;%d;%d;%dm", red, green, blue)
	b2 := []byte(s)
	return append(b[:], b2...)
}

func exportClassicANSI(image *reximage.RexImage, outFile *os.File) {

	// Draw it back out
	for i := 0; i < int(image.LayerCount); i++ {
		layer := &image.Layers[i]

		strideRemaining := int(layer.Width)

		lastFG := reximage.RexRGB{0, 0, 0}
		fgReset := true

		lastBG := reximage.RexRGB{0, 0, 0}
		bgReset := true

		for j := 0; j < int(layer.Height*layer.Width); j++ {
			cell := &layer.Cells[j]

			glyph := cell.Glyph

			if cell.IsTransparent() {
				outFile.Write(classicReset())
				outFile.Write([]byte{' '})
				fgReset = true
				bgReset = true
			} else {
				// color
				if fgReset || !reximage.CompareRGB(cell.Fg, lastFG) {
					outFile.Write(classicFg24(cell.Fg.Red, cell.Fg.Green, cell.Fg.Blue))
					lastFG = cell.Fg
					fgReset = false
				}

				if bgReset || !reximage.CompareRGB(cell.Bg, lastBG) {
					outFile.Write(classicBg24(cell.Bg.Red, cell.Bg.Green, cell.Bg.Blue))
					lastBG = cell.Bg
					bgReset = false
				}

				// glyph
				outFile.Write([]byte{glyph})
			}

			strideRemaining--

			if strideRemaining == 0 {
				strideRemaining = int(layer.Width)
				outFile.Write(classicReset())
				outFile.Write([]byte{'\n'})
				fgReset = true
				bgReset = true
				continue

			}
		}
	}
}
