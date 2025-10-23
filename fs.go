package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func getRandomLetter() string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randomIndex := rand.Intn(len(alphabet))
	return string(alphabet[randomIndex])
}

func AbsPath(p string) (string, error) {
	absolutePath, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}
	return absolutePath, nil
}

func TestFilesExist(s Settings) bool {
	for i := range s.Streams {
		fileName := fmt.Sprintf("%s/testfile_%d.txt", s.SourceDir, i)
		fi, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			return false
		}
		if err != nil {
			return false
		}
		if fi.Size() != int64(s.FileSizeBytes) {
			return false
		}
	}
	return true
}

func CreateTestFiles(s Settings) error {
	errChan := make(chan error)
	var wg sync.WaitGroup

	for i := range s.Streams {
		wg.Go(func() {
			fileName := fmt.Sprintf("%s/testfile_%d.txt", s.SourceDir, i)
			file, err := os.Create(fileName)
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()

			writer := bufio.NewWriter(file)
			defer writer.Flush()

			for range s.FileSizeBytes {
				if _, err := writer.Write([]byte(getRandomLetter())); err != nil {
					errChan <- err
					return
				}
			}
		})
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		return err
	}

	return nil
}

func StartTransfers(s Settings) error {
	errChan := make(chan error)
	cancelChan := make(chan struct{})
	var wg sync.WaitGroup

	for i := range s.Streams {
		wg.Go(func() {
			startTime := time.Now()
			select {
			case <-cancelChan:
				return
			default:
				srcFile := fmt.Sprintf("%s/testfile_%d.txt", s.SourceDir, i)
				dstFile := fmt.Sprintf("%s/testfile_%d.txt", s.DestDir, i)

				src, err := os.Open(srcFile)
				if err != nil {
					errChan <- fmt.Errorf("error opening source file %s: %v", srcFile, err)
					cancelChan <- struct{}{}
					return
				}
				defer src.Close()

				dst, err := os.Create(dstFile)
				if err != nil {
					errChan <- fmt.Errorf("error creating destination file %s: %v", dstFile, err)
					cancelChan <- struct{}{}
					return
				}
				defer dst.Close()

				if _, err := io.Copy(dst, src); err != nil {
					errChan <- fmt.Errorf("error transferring file %s to %s: %v", srcFile, dstFile, err)
					cancelChan <- struct{}{}
					return
				}
				duration := time.Since(startTime).Truncate(time.Millisecond)
				avgSpeed := float64(s.FileSizeGB) / duration.Seconds()

				fmt.Printf("%d: %s @ %.2f GB/s\n", i, duration.String(), avgSpeed)
			}
		})
	}
	wg.Wait()
	close(cancelChan)
	close(errChan)
	for err := range errChan {
		return err
	}

	return nil
}

func CleanupTestFiles(s Settings) []error {
	var errs []error

	for i := range s.Streams {

		srcFile := fmt.Sprintf("%s/testfile_%d.txt", s.SourceDir, i)
		dstFile := fmt.Sprintf("%s/testfile_%d.txt", s.DestDir, i)

		if err := os.Remove(srcFile); err != nil {
			errs = append(errs, fmt.Errorf("error removing source file %s: %v", srcFile, err))
		}

		if err := os.Remove(dstFile); err != nil {
			errs = append(errs, fmt.Errorf("error removing destination file %s: %v", dstFile, err))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
