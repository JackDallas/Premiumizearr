package directory_watcher

// WatchDirectory watches a directory for changes.
type WatchDirectory struct {
	// Path is the path to the directory to watch.
	Path string
	// Filter is the filter to apply to the directory.
	Filter string
	// Recursive is true if the directory should be watched recursively.
	Recursive bool
	// MatchFunction is the function to use to match files.
	MatchFunction func(string) bool
	// Callback is the function to call when a file is created that matches with MatchFunction.
	CallbackFunction func(string)
}
