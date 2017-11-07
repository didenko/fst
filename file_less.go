// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"

import (
	"os"
)

// FileLess describes functions which can be provided to
// compare two os.FileInfo structs and related files
type FileLess func(left, right os.FileInfo) bool

// LessName is basic for comparing directories and should
// be provided as a first comparator in most cases
func LessName(left, right os.FileInfo) bool {
	return left.Name() < right.Name()
}

// LessDir differentiates directorries vs. files and puts
// directories earlier in a sort order
func LessDir(left, right os.FileInfo) bool {
	return left.IsDir() && !right.IsDir()
}

// LessSize compares sizes of files, given that both of the
// files are regular files as opposed to not directories, etc.
func LessSize(left, right os.FileInfo) bool {
	return left.Mode().IsRegular() &&
		right.Mode().IsRegular() &&
		(left.Size() < right.Size())
}

// LessTime compares files' last modification times
func LessTime(left, right os.FileInfo) bool {
	return left.ModTime().Before(right.ModTime())
}

// LessPerm compares bits 0-8 of Unix-like file permissions
func LessPerm(left, right os.FileInfo) bool {
	return left.Mode().Perm() < right.Mode().Perm()
}

// Less applies provided comparators to the pair of os.FileInfo structs.
func Less(left, right os.FileInfo, cmps ...FileLess) bool {
	for _, less := range cmps {
		if less(left, right) {
			return true
		}
	}
	return false
}
