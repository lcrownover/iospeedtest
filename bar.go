package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func StartTransferBars(s Settings) float64 {
	var wg sync.WaitGroup
	p := mpb.New(
		mpb.WithWaitGroup(&wg),
		mpb.WithWidth(64),
		mpb.WithAutoRefresh(),
	)

	sumGiBps := 0.0
	perStreamGiBps := make([]float64, s.Streams)
	var mu sync.Mutex

	for i := range s.Streams {
		// Create destination file
		outputFileName := fmt.Sprintf("%s/iospeedtest_%d.txt", s.DestDir, i)
		f, err := os.Create(outputFileName)
		if err != nil {
			fmt.Printf("error creating file %s: %v", outputFileName, err)
			return 0.0
		}
		defer f.Close()

		// Set up the bar
		stream := fmt.Sprintf("Stream#%02d:", i)
		bar := p.AddBar(s.FileSizeBytes,
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(stream, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
				decor.Counters(decor.SizeB1024(0), "% .2f / % .2f ", decor.WCSyncWidth),
				decor.OnComplete(
					decor.Name("transferring", decor.WCSyncSpaceR),
					"done!",
				),
			),
			mpb.AppendDecorators(
				decor.AverageSpeed(decor.SizeB1024(0), "% .2f", decor.WCSyncWidth),
				decor.OnComplete(decor.Percentage(decor.WC{W: 5}), ""),
			),
		)

		// create a proxy reader
		proxyReader := bar.ProxyReader(io.LimitReader(rand.Reader, s.FileSizeBytes))
		defer func() {
			proxyReader.Close()
		}()

		wg.Go(func() {
			start := time.Now()
			writtenBytes, err := io.Copy(f, proxyReader)
			elapsed := time.Since(start).Seconds()
			if elapsed > 0 {
				speed := float64(writtenBytes) / elapsed
				mu.Lock()
				perStreamGiBps[i] = speed / (1024 * 1024 * 1024)
				mu.Unlock()
			}

			if err != nil {
				fmt.Printf("error copying to file %s: %v\n", outputFileName, err)
				return
			}
		})

	}

	p.Wait()

	for _, v := range perStreamGiBps {
		sumGiBps += v
	}

	return sumGiBps
}
