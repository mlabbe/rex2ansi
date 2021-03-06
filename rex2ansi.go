package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/mlabbe/rex2ansi/reximage"
)

var (
	verbose     = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	skipFlatten = kingpin.Flag("skip-flatten", "Don't flatten image").Short('s').Bool()
	onlyUTF8    = kingpin.Flag("only-utf8", "Only generate utf-8 ANSI").Bool()
	onlyCP437   = kingpin.Flag("only-cp437", "Only codepage 437 (classic) ANSI").Bool()
	outputDir   = kingpin.Flag("output-dir", "Directory to write files to").Short('o').Default(".").ExistingDir()
	version     = kingpin.Flag("version", "print version and exit").Bool()

	// positional, bash wildcard-friendly
	paths = kingpin.Arg("files", "files to operate on").ExistingFiles()

	Version = 2
)

func getOutPath(inFile string, utf8 bool) string {
	// get the filename without extension or path
	baseName := strings.TrimSuffix(inFile, filepath.Ext(inFile))
	baseName = filepath.Base(baseName)

	ext := "ans"
	if utf8 {
		ext = "u8ans" /* I made this up */
	}

	return *outputDir + "/" + baseName + "." + ext
}

func main() {

	kingpin.Parse()

	if *version {
		log.Printf("rex2ansi version %d", Version)
		os.Exit(0)
	}

	if len(*paths) == 0 {
		log.Printf("No files specified.  Try --help")
		os.Exit(1)
	}

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
			continue
		}

		if !*skipFlatten {
			image.Flatten()
		}

		//
		// Write output files
		//

		// utf-8
		if !*onlyCP437 {
			outPath := getOutPath(path, true)
			if *verbose {
				log.Printf("Writing file: %s", outPath)
			}

			outH, err := os.Create(outPath)
			if err != nil {
				log.Fatal(err)
			}
			defer outH.Close()

			exportUTF8ANSI(image, outH)
		}

		// cp437
		if !*onlyUTF8 {
			outPath := getOutPath(path, false)
			if *verbose {
				log.Printf("Writing file: %s", outPath)
			}

			outH, err := os.Create(outPath)
			if err != nil {
				log.Fatal(err)
			}
			defer outH.Close()

			exportClassicANSI(image, outH)
		}
	}

	os.Exit(errorCount)
}
