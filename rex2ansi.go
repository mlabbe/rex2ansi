package main

//
// todos:
//
// 8. unflattened files should use cursor manip
// 9. gofmt
// 10. switch from convert-utf 8 to no-utf8, do both by default

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
	outputDir   = kingpin.Flag("output-dir", "Directory to write files to").Short('o').Default(".").ExistingDir()

	// positional, bash wildcard-friendly
	paths   = kingpin.Arg("files", "files to operate on").Required().ExistingFiles()
)

func getOutPath(inFile string) string {
	// get the filename without extension or path
	baseName := strings.TrimSuffix(inFile, filepath.Ext(inFile))
	baseName = filepath.Base(baseName)

	ext := "ans"
	if *convUTF8 {
		ext = "u8ans" /* I made this up */
	}

	return *outputDir + "/" + baseName + "." + ext
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
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", path, err)
			continue;
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
