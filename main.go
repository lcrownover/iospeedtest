package main

import (
	"flag"
	"fmt"
	"os"
)

type Settings struct {
	SourceDir     string
	DestDir       string
	FileSizeBytes int64
	FileSizeGB    int64
	Streams       int
	Cleanup       bool
}

func NewSettings(srcDir, dstDir string, size int64, streams int, cleanup bool) Settings {
	return Settings{
		SourceDir:     srcDir,
		DestDir:       dstDir,
		FileSizeBytes: size * 1024 * 1024 * 1024,
		FileSizeGB:    size,
		Streams:       streams,
		Cleanup:       cleanup,
	}
}

func main() {
	srcDir := flag.String("srcdir", "", "source directory")
	dstDir := flag.String("dstdir", "", "destination directory")
	size := flag.Int64("size", 1, "file size in GB")
	streams := flag.Int("streams", 1, "number of streams")
	cleanup := flag.Bool("cleanup", false, "cleanup all test files when finished")
	flag.Parse()

	if *srcDir == "" || *dstDir == "" {
		fmt.Println("source directory and destination directory are required")
		os.Exit(1)
	}
	srcFullPath, err := AbsPath(*srcDir)
	if err != nil {
		fmt.Println("error getting absolute path for source directory:", err)
		os.Exit(1)
	}
	dstFullPath, err := AbsPath(*dstDir)
	if err != nil {
		fmt.Println("error getting absolute path for destination directory:", err)
		os.Exit(1)
	}

	settings := NewSettings(srcFullPath, dstFullPath, *size, *streams, *cleanup)

	fmt.Println("Source directory:\t", settings.SourceDir)
	fmt.Println("Destination directory:\t", settings.DestDir)
	fmt.Println("File size (bytes):\t", settings.FileSizeBytes)
	fmt.Println("File size (GB):\t\t", settings.FileSizeGB)
	fmt.Println("Number of streams:\t", settings.Streams)

	fmt.Println("Creating randomized source files...")

	if !TestFilesExist(settings) {
		err := CreateTestFiles(settings)
		if err != nil {
			fmt.Println("error creating test files:", err)
			os.Exit(1)
		}
	}

	fmt.Println("Starting transfer(s)...")
	StartTransfers(settings)

	if settings.Cleanup {
		fmt.Println("Cleaning up test files...")
		errs := CleanupTestFiles(settings)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Println("error cleaning up test files:", err)
			}
			os.Exit(1)
		}
	}
}
