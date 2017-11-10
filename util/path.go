package util

import (
	"os"
	"os/user"
	"path/filepath"
)

// UserDir returns path to user home directory
func UserDir() string {
	currentUser, _ := user.Current()
	return currentUser.HomeDir
}

// ResolvePath makes absolute path from a path, relative to another
// Returns the absolute path and a boolean value which indicates
// if the returned path exists or not
func ResolvePath(path string, relpath string) (string, bool) {

	// Deal with ~/ form - substitute user's home path
	if filepath.HasPrefix(path, "~/") {
		path = filepath.Join(UserDir(), path[2:])
	}

	// Join relpath with path
	if filepath.IsAbs(path) == false {
		path = filepath.Join(relpath, path)
	}

	// Clean up the path
	path = filepath.Clean(path)

	// Determine if path exists
	exists := true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exists = false
	}

	// Return
	return path, exists
}
