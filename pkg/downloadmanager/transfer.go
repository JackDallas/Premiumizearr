package downloadmanager

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

var (
	ErrorNoTransferWithID = errors.New("no transfer with id")
)

type transferStatus int

const (
	STATUS_QUEUED transferStatus = iota
	STATUS_DOWNLOADING
	STATUS_PAUSED
	STATUS_COMPLETED
	STATUS_CANCELED
	STATUS_ERROR
)

type Transfer struct {
	id               int64
	totalSize        atomic.Int64
	downloaded       atomic.Int64
	savePath         string
	url              string
	urlLock          sync.Mutex
	status           transferStatus
	statusLock       sync.Mutex
	errorStrings     []string
	errorStringsLock sync.Mutex
	tempFileName     string
	Finished         chan bool
}

func NewTransfer(id int64, url string, savePath string) Transfer {
	return Transfer{
		id:               id,
		totalSize:        atomic.Int64{},
		downloaded:       atomic.Int64{},
		savePath:         savePath,
		url:              url,
		urlLock:          sync.Mutex{},
		status:           STATUS_QUEUED,
		statusLock:       sync.Mutex{},
		errorStrings:     make([]string, 0),
		errorStringsLock: sync.Mutex{},
		tempFileName:     "",
	}
}

func (t *Transfer) SetID(id int64) {
	atomic.StoreInt64(&t.id, id)
}

func (t *Transfer) GetID() int64 {
	return atomic.LoadInt64(&t.id)
}

func (t *Transfer) SetTotalSize(size int64) {
	t.totalSize.Store(size)
}

func (t *Transfer) GetTotalSize() int64 {
	return t.totalSize.Load()
}

func (t *Transfer) SetDownloaded(size int64) {
	t.downloaded.Store(size)
}

func (t *Transfer) GetDownloaded() int64 {
	return t.downloaded.Load()
}

func (t *Transfer) SetURL(url string) {
	t.urlLock.Lock()
	t.url = url
	t.urlLock.Unlock()
}

func (t *Transfer) GetURL() string {
	t.urlLock.Lock()
	defer t.urlLock.Unlock()
	return t.url
}

func (t *Transfer) SetStatus(status transferStatus) {
	t.statusLock.Lock()
	t.status = status
	t.statusLock.Unlock()
}

func (t *Transfer) GetStatus() transferStatus {
	t.statusLock.Lock()
	defer t.statusLock.Unlock()
	return t.status
}

func (t *Transfer) AddErrorString(str string) {
	t.errorStringsLock.Lock()
	t.errorStrings = append(t.errorStrings, str)
	t.errorStringsLock.Unlock()
}

func (t *Transfer) GetErrorStrings() []string {
	t.errorStringsLock.Lock()
	defer t.errorStringsLock.Unlock()
	return t.errorStrings
}

func (t *Transfer) GetTempFilePath() string {
	if t.tempFileName == "" {
		url, err := url.Parse(t.GetURL())
		if err != nil {
			t.tempFileName = fmt.Sprintf("download-%d", t.GetID())
		} else {
			finalPath := strings.Split(url.Path, "/")[len(strings.Split(url.Path, "/"))-1]
			t.tempFileName = fmt.Sprintf("download-%d-%s", t.GetID(), finalPath)
		}
	}

	return t.tempFileName
}

func (t *Transfer) Write(p []byte) (int, error) {
	if t.GetStatus() == STATUS_CANCELED || t.GetStatus() == STATUS_PAUSED {
		return 0, io.EOF
	}
	t.SetDownloaded(t.GetDownloaded() + int64(len(p)))
	return len(p), nil
}

func (t *Transfer) Pause() error {
	t.SetStatus(STATUS_PAUSED)
	return nil
}

func (t *Transfer) Cancel() error {
	t.SetStatus(STATUS_CANCELED)
	t.Finished <- true
	return nil
}

func (t *Transfer) Resume() error {
	return t.Download()
}

func (t *Transfer) Download() error {
	client := &http.Client{}

	//Built http get request with a content range header
	req, err := http.NewRequest("GET", t.GetURL(), nil)
	if err != nil {
		return err
	}
	if t.GetDownloaded() > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", t.GetDownloaded()))
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	t.SetTotalSize(resp.ContentLength)
	var out *os.File

	if t.GetDownloaded() > 0 {
		out, err = os.Open(t.GetTempFilePath())
	} else {
		out, err = os.Create(t.GetTempFilePath())
	}

	if err != nil {
		return err
	}

	t.SetStatus(STATUS_DOWNLOADING)
	go func() {
		defer out.Close()

		if _, err := io.Copy(out, io.TeeReader(resp.Body, t)); err != nil {
			t.AddErrorString(err.Error())
			t.SetStatus(STATUS_ERROR)
			log.Error(err)
		}
		t.SetStatus(STATUS_COMPLETED)
		t.Finished <- true
	}()

	return nil
}

func Start() {

}
