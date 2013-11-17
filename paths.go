/*
Package paths provides helpers for dealing with files and directories.
*/
package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type VCSKind int

// The various types of vcs directories that FindVCSRoot can find.
const (
	NONE VCSKind = iota
	GIT
	HG
	BZR
	SVN
)

const (
	_PathSeparator = string(os.PathSeparator)
)

var vcsKinds = []VCSKind{GIT, HG, BZR, SVN}
var vcsDirs = []string{".git", ".hg", ".bzr", ".svn"}

// FindVCSRoot finds a .git/.hg/.bzr/.svn by stepping up the given path.
// returns an empty string if nothing could be found.
func FindVCSRoot(path string) (VCSKind, string, error) {
	for ; len(path) != 0; path = WalkUpPath(path) {
		for i, vcsDir := range vcsDirs {
			yes, err := DirExists(filepath.Join(path, vcsDir))
			if err != nil {
				return NONE, "", err
			} else if yes {
				return vcsKinds[i], path, nil
			}
		}
	}

	return NONE, "", nil
}

// WalkUpPath removes the head of path.
// /home/user -> /home
// /home/user/file.txt -> /home/user
// If path is empty or is root returns "".
// If path has no os.PathSeperator in it, returns "".
func WalkUpPath(path string) string {
	if len(path) == 0 || path == _PathSeparator {
		return ""
	}
	path = filepath.Clean(path)
	index := strings.LastIndex(path, _PathSeparator)
	if index == 0 {
		return _PathSeparator
	} else if index < 0 {
		return ""
	}
	return path[:index]
}

// EnsureDirectory ensures a directory exists, or it creates it. Returns
// true if the directory had to be created.
func EnsureDirectory(dir string) (bool, error) {
	if exists, err := DirExists(dir); err != nil {
		return false, err
	} else if exists {
		return false, nil
	}

	err := os.MkdirAll(dir, 0770)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DirExists checks to see if a directory exists.
func DirExists(dir string) (bool, error) {
	f, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return false, err
	}
	if !f.IsDir() {
		return false, fmt.Errorf("Expected %s to be dir, but found file.", dir)
	}
	return true, nil
}

// FileExists checks to see if a directory exists.
func FileExists(file string) (bool, error) {
	f, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return false, err
	}
	if f.IsDir() {
		return false, fmt.Errorf("Expected %s to be file, but found dir.", file)
	}
	return true, nil
}
