package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type CLI struct {
	DstDir  string `arg:"" name:"destination" short:"d" help:"Destination directory." type:"existingdir"`
	Size    int64  `name:"gb" help:"How many gigabytes to shovel." default:"1"`
	Streams int    `name:"streams" short:"s" help:"Number of streams." default:"1"`
	Cleanup bool   `name:"cleanup" short:"c" help:"Cleanup destination files when finished."`
}

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
	var cli CLI
	kong.Parse(&cli, kong.Name("iospeedtest"), kong.Description("A simple tool to measure how fast bits fly to the destination directory."), kong.UsageOnError())

	dstFullPath, err := AbsPath(cli.DstDir)
	if err != nil {
		fmt.Println("error getting absolute path for destination directory:", err)
		os.Exit(1)
	}
	err = CheckDestDirExists(dstFullPath)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	settings := NewSettings(dstFullPath, cli.Size, cli.Streams, cli.Cleanup)

	fmt.Println("Destination directory:\t", settings.DestDir)
	fmt.Println("Number of streams:\t", settings.Streams)
	fmt.Println("Bytes per stream:\t", settings.FileSizeBytes)

	sumSpeed := StartTransferBars(settings)
	fmt.Println()
	fmt.Println("Total data transferred:\t", fmt.Sprintf("%d GiB", settings.FileSizeGB*int64(settings.Streams)))
	fmt.Println("Aggregate speed:\t", fmt.Sprintf("%.2f GiB/s", sumSpeed))
	fmt.Println("Avg speed per stream:\t", fmt.Sprintf("%.2f GiB/s", sumSpeed/float64(settings.Streams)))

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
