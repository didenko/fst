// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"bufio"
	"os"
	"testing"
)

// FileRank is the signature of functions which can be
// provided to TreeDiff to compare two *FileInfoPath structs
// and related files
type FileRank func(left, right *FileInfoPath) bool

// ByName is basic for comparing directories and should
// be provided as a first comparator in most cases
func ByName(left, right *FileInfoPath) bool {
	return left.Name() < right.Name()
}

// ByDir differentiates directorries vs. files and puts
// directories earlier in a sort order
func ByDir(left, right *FileInfoPath) bool {
	return left.IsDir() && !right.IsDir()
}

// BySize compares sizes of files, given that both of the
// files are regular files as opposed to not directories, etc.
func BySize(left, right *FileInfoPath) bool {
	return left.Mode().IsRegular() &&
		right.Mode().IsRegular() &&
		(left.Size() < right.Size())
}

// ByTime compares files' last modification times
func ByTime(left, right *FileInfoPath) bool {
	return left.ModTime().Before(right.ModTime())
}

// ByPerm compares bits 0-8 of Unix-like file permissions
func ByPerm(left, right *FileInfoPath) bool {
	return left.Mode().Perm() < right.Mode().Perm()
}

// ByContent returns a function which compares files'
// content without first comparing sizes. For example,
// file containing "aaa" will rank as lesser than the one
// containing "ab" even though it is opposite to their sizes.
// To consider sizes first, make sure to specify the BySize
// comparator earlier in the chain.
func ByContent(t *testing.T) FileRank {
	return func(left, right *FileInfoPath) bool {
		leftF, err := os.Open(left.Path())
		if err != nil {
			t.Fatal(err)
		}
		defer leftF.Close()

		rightF, err := os.Open(right.Path())
		if err != nil {
			t.Fatal(err)
		}
		defer rightF.Close()

		leftBuf := bufio.NewReader(leftF)
		rightBuf := bufio.NewReader(rightF)

		for {
			rByte, err := rightBuf.ReadByte()
			if err != nil {
				return false
			}

			lByte, err := leftBuf.ReadByte()
			if err != nil {
				return true
			}

			if lByte == rByte {
				continue
			}

			if lByte < rByte {
				return true
			}

			return false
		}
	}
}

// Less applies provided comparators to the pair of *FileInfoPath structs.
func Less(left, right *FileInfoPath, cmps ...FileRank) bool {
	for _, less := range cmps {
		if less(left, right) {
			return true
		}
	}
	return false
}
