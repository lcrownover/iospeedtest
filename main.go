package main

import (
	"flag"
	"fmt"
	"os"
)

type Settings struct {
	DestDir       string
	FileSizeBytes int64
	FileSizeGB    int64
	Streams       int
	Cleanup       bool
}

func NewSettings(dstDir string, size int64, streams int, cleanup bool) Settings {
	return Settings{
		DestDir:       dstDir,
		FileSizeBytes: size * 1024 * 1024 * 1024,
		FileSizeGB:    size,
		Streams:       streams,
		Cleanup:       cleanup,
	}
}

// func main() {
// 	run()
// }

func main() {
	dstDir := flag.String("dstdir", "", "destination directory")
	size := flag.Int64("size", 1, "file size in GB")
	streams := flag.Int("streams", 1, "number of streams")
	cleanup := flag.Bool("cleanup", false, "cleanup all test files when finished")
	flag.Parse()

	if *dstDir == "" {
		fmt.Println("destination directory is required")
		os.Exit(1)
	}
	dstFullPath, err := AbsPath(*dstDir)
	if err != nil {
		fmt.Println("error getting absolute path for destination directory:", err)
		os.Exit(1)
	}
	err = CheckDestDirExists(dstFullPath)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	settings := NewSettings(dstFullPath, *size, *streams, *cleanup)

	fmt.Println("Destination directory:\t", settings.DestDir)
	fmt.Println("Number of streams:\t", settings.Streams)
	fmt.Println("Bytes per stream:\t", settings.FileSizeBytes)

	StartTransferBars(settings)

	if settings.Cleanup {
		errs := CleanupTestFiles(settings)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Println("error cleaning up test files:", err)
			}
			os.Exit(1)
		}
	}
}
