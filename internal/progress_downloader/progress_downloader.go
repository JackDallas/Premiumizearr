package progress_downloader

// https://golangcode.com/download-a-file-with-progress/

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	StartTime  time.Time
	LastUpdate time.Time
	LastAmount uint64
	Total      uint64
}

func NewWriteCounter() *WriteCounter {
	return &WriteCounter{
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
		LastAmount: 0,
		Total:      0,
	}
}

func (wc *WriteCounter) GetSpeed() string {
	timeSince := time.Since(wc.StartTime).Seconds()

	if timeSince == 0 || timeSince < 1 {
		return "0 B/s"
	}

	return fmt.Sprintf("%s/s", humanize.Bytes(wc.Total/uint64(timeSince)))
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.LastAmount = wc.Total
	wc.Total += uint64(n)
	return n, nil
}

func (wc WriteCounter) GetProgress() string {
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	return fmt.Sprintf("%s complete", humanize.Bytes(wc.Total))
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(url string, filepath string, counter *WriteCounter) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}
