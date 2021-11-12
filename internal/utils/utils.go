package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func StripDownloadTypesExtention(fileName string) string {
	var exts = [...]string{".nzb", ".magnet"}
	for _, ext := range exts {
		fileName = strings.TrimSuffix(fileName, ext)
	}

	return fileName
}

func GetTempBaseDir() string {
	return path.Join(os.TempDir(), "premiumizearrd")
}

func GetTempDir() (string, error) {
	// Create temp dir in os temp location
	tempDir := GetTempBaseDir()
	err := os.Mkdir(tempDir, os.ModePerm)
	dir, err := ioutil.TempDir(tempDir, "unzip-")
	if err != nil {
		return "", err
	}
	return dir, nil
}

func DownloadAndExtractZip(link string, destination string) error {
	tempDir, err := GetTempDir()
	if err != nil {
		return err
	}

	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP status code: %d When trying to download %s", resp.StatusCode, link)
	}
	splitString := strings.Split(link, "/")
	savePath := path.Join(tempDir, splitString[len(splitString)-1])

	log.Trace("Downloading to: ", savePath)
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Tracef("Unzipping %s to %s", savePath, destination)
	err = Unzip(savePath, destination)
	if err != nil {
		return err
	}

	log.Tracef("Removing zip %s from system", savePath)
	err = os.RemoveAll(savePath)
	if err != nil {
		return err
	}

	return nil
}

// https://golangcode.com/unzip-files-in-go/
func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: https://snyk.io/research/zip-slip-vulnerability#go
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func StringInSlice(a string, list []string) int {
	for i, b := range list {
		if b == a {
			return i
		}
	}
	return -1
}
