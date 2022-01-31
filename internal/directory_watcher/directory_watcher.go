package directory_watcher

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// NewWatchDirectory creates a new WatchDirectory.
func NewDirectoryWatcher(path string, recursive bool, matchFunction func(string) bool, callbackFunction func(string)) *WatchDirectory {
	return &WatchDirectory{
		Path: path,
		// TODO (Unused): Add recursive abilities
		Recursive:        recursive,
		MatchFunction:    matchFunction,
		CallbackFunction: callbackFunction,
	}
}

func (w *WatchDirectory) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					if w.MatchFunction(event.Name) {
						w.CallbackFunction(event.Name)
					}
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	cleanPath := filepath.Clean(w.Path)
	_, err = os.Stat(cleanPath)
	if os.IsNotExist(err) {
		return err
	}

	err = watcher.Add(cleanPath)
	if err != nil {
		return err
	}
	<-done
	return nil
}
