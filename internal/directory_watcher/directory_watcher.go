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
	var err error
	w.Watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-w.Watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					if w.MatchFunction(event.Name) {
						w.CallbackFunction(event.Name)
					}
				}
			case _, ok := <-w.Watcher.Errors:
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

	err = w.Watcher.Add(cleanPath)
	if err != nil {
		return err
	}

	return nil
}

func (w *WatchDirectory) UpdatePath(path string) error {
	w.Watcher.Remove(w.Path)
	w.Path = path
	return w.Watcher.Add(w.Path)
}

func (w *WatchDirectory) Stop() error {
	return w.Watcher.Close()
}
