package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jackdallas/premiumizearr/internal/utils"
	"github.com/jackdallas/premiumizearr/pkg/directorywatcher"
	"github.com/jackdallas/premiumizearr/pkg/stringqueue"
	log "github.com/sirupsen/logrus"
)

type DirectoryWatcherService struct {
	premiumizearrd *premiumizearrd
	Queue          *stringqueue.StringQueue
}

//TODO (Radarr): accept paths as a parameter, support multiple paths
//Watch: This is the entrypoint for the directory watcher
func (dw *DirectoryWatcherService) Watch() {
	log.Info("Starting directory watcher...")

	log.Info("Clearing tmp directory...")
	tempDir := utils.GetTempBaseDir()
	err := os.RemoveAll(tempDir)
	if err != nil {
		log.Errorf("Error clearing tmp directory %s", tempDir)
	}
	os.Mkdir(tempDir, os.ModePerm)

	log.Info("Creating Queue...")
	dw.Queue = stringqueue.NewStringQueue()

	log.Info("Starting uploads processor...")
	go dw.processUploads()

	log.Info("Starting initial directory scan...")
	go dw.initialDirectoryScan()

	// Build and start a DirectoryWatcher
	watcher := directorywatcher.NewDirectoryWatcher(dw.premiumizearrd.Config.SonarrBlackholeDirectory,
		false,
		dw.checkFile,
		dw.addFileToQueue,
	)

	// Block
	watcher.Watch()
}

func (dw *DirectoryWatcherService) initialDirectoryScan() {
	log.Trace("Initial directory scan")
	files, err := ioutil.ReadDir(dw.premiumizearrd.Config.SonarrBlackholeDirectory)
	if err != nil {
		log.Errorf("Error with initial directory scan %+v", err)
	}

	for _, file := range files {
		go func(file os.FileInfo) {
			file_path := path.Join(dw.premiumizearrd.Config.SonarrBlackholeDirectory, file.Name())
			if dw.checkFile(file_path) {
				dw.addFileToQueue(file_path)
			}
		}(file)
	}
}

func (dw *DirectoryWatcherService) checkFile(path string) bool {
	log.Tracef("Checking file %s", path)

	fi, err := os.Stat(path)
	if err != nil {
		log.Errorf("Error checking file %s", path)
		return false
	}

	if fi.IsDir() {
		log.Errorf("Directory created in blackhole %s ignoring (Warning premiumizearrzed does not look in subfolders!)", path)
		return false
	}

	ext := filepath.Ext(path)
	if ext == ".nzb" || ext == ".magnet" {
		return true
	} else {
		return false
	}
}

func (dw *DirectoryWatcherService) addFileToQueue(path string) {
	dw.Queue.Add(path)
	log.Infof("File created in blackhole %s added to queue. Queue length %d", path, dw.Queue.Len())
}

func (dw *DirectoryWatcherService) processUploads() {
	//TODO: Global running state
	for {
		if dw.Queue.Len() < 1 {
			log.Trace("No files in queue, sleeping for 10 seconds")
			time.Sleep(time.Second * time.Duration(10))
		}
		sleepTimeSeconds := 2
		filePath := dw.Queue.GetTopOfQueue()

		if filePath != "" {
			log.Debugf("Processing %s", filePath)
			err := dw.premiumizearrd.premiumizearrClient.CreateTransfer(filePath)
			if err != nil {
				if err.Error() == "Limit of transfers reached!" {
					log.Info("Transfer limit reached waiting 10 seconds and retrying ")
					sleepTimeSeconds = 10
				} else {
					log.Error(err)
				}
			} else {
				os.Remove(filePath)
				if err != nil {
					log.Errorf("Error could not delete %s Error: %+v", filePath, err)
				}
				dw.Queue.DeleteTopOfQueue()
				log.Infof("Removed %s from blackhole queue. Queue Size: %d", filePath, dw.Queue.Len())
			}
			time.Sleep(time.Second * time.Duration(sleepTimeSeconds))
		} else {
			// Empty link
			dw.Queue.DeleteTopOfQueue()
			log.Errorf("Removed %s from blackhole queue as it's an empty path. Queue Size: ", filePath, dw.Queue.Len())
		}
	}
}
