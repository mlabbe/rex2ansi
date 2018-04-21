package main

//
// todos:
//
// 4. Dont' fatal out in Read(), return error
// 5. Support --output-dir
// 6. Handle stdio
// 7. Error out on input file not being rexpaint.
// 8. unflattened files should use cursor manip
// 9. gofmt

import (
	"os"
	"log"
	"fmt"
	"path/filepath"
	"strings"
	"github.com/alecthomas/kingpin"
	"frogtoss.com/rex2ansi/reximage"
)

var (
	verbose     = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	skipFlatten = kingpin.Flag("skip-flatten", "Don't flatten image").Short('s').Bool()
	convUTF8    = kingpin.Flag("convert-utf8", "Convert codepage 437 to utf-8").Short('c').Bool()

	// positional, bash wildcard-friendly
	paths   = kingpin.Arg("files", "files to operate on").Required().ExistingFiles()
)

func getOutPath(inFile string) string {
	baseName := strings.TrimSuffix(inFile, filepath.Ext(inFile))

	ext := "ans"
	if *convUTF8 {
		ext = "u8ans" /* I made this up */
	}

	return baseName + "." + ext
}

func main() {

	kingpin.Parse()

	errorCount := 0
	for _, path := range *paths {
		if *verbose {
			log.Printf("Reading File: %s", path)
		}

		//
		// Read Rexpaint File
		//
		inFile, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer inFile.Close()

		image, err := reximage.Read(inFile, *verbose)
		if err != nil {
			errorCount++
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		if !*skipFlatten {
			image.Flatten()
		}

		//
		// Write classic ansi file
		//
		outPath := getOutPath(path)
		if *verbose {
			log.Printf("Writing File: %s", outPath)
		}
		outFile, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()

		if *convUTF8 {
			exportUTF8ANSI(image, outFile)
		} else {
			exportClassicANSI(image, outFile)
		}
	}

	os.Exit(errorCount)
}
