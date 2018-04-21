package main

import (
	"os"
	"fmt"
	"frogtoss.com/rex2ansi/reximage"
)

// ansiclassic is the catch-all phrase for codepage 437 "ASCII" ansi,
// which is what was used on DOS.  In order to view these files you
// must have codepage 437 installed on your terminal with a font that
// has the codepage 437 charset.  This is not modern, but it is what
// you will need to display these art files in a Windows console.


func reset() string {
	return "\u001b[0m"
}

func fg24(red, green, blue byte) string {
	return fmt.Sprintf("\u001b[38;2;%d;%d;%dm", red, green, blue)
}

func bg24(red, green, blue byte) string {
	return fmt.Sprintf("\u001b[48;2;%d;%d;%dm", red, green, blue)
}

func exportUTF8ANSI(image *reximage.RexImage, outFile *os.File) {

	// Draw it back out
	for i := 0; i < int(image.LayerCount); i++ {
		layer := &image.Layers[i]

		strideRemaining := int(layer.Width)

		lastFG := reximage.RexRGB{0,0,0}
		fgReset := true

		lastBG := reximage.RexRGB{0,0,0}
		bgReset := true

		for j := 0; j < int(layer.Height*layer.Width); j++ {
			cell := &layer.Cells[j]

			glyph := cell.Glyph

			if cell.IsTransparent() {
				outFile.WriteString(reset() + " ")
				fgReset = true
				bgReset = true
			} else {
				// color
				if fgReset || !reximage.CompareRGB(cell.Fg, lastFG) {
					outFile.WriteString(fg24(cell.Fg.Red, cell.Fg.Green, cell.Fg.Blue))
					lastFG = cell.Fg
					fgReset = false
				}

				if bgReset || !reximage.CompareRGB(cell.Bg, lastBG) {
					outFile.WriteString(bg24(cell.Bg.Red, cell.Bg.Green, cell.Bg.Blue))
					lastBG = cell.Bg
					bgReset = false
				}

				// glyph
				// implicit conversion to utf-8
				outFile.WriteString(fmt.Sprintf("%c", glyph))
			}

			strideRemaining--

			if strideRemaining == 0 {
				strideRemaining = int(layer.Width)
				outFile.WriteString(reset() + "\n")
				fgReset = true
				bgReset = true
				continue
			}
		}
	}

}
