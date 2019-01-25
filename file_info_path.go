// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"os"
)

// FileInfoPath is a wrapper of os.FileInfo with an additional
// field to store the path to the file of interest
type FileInfoPath struct {
	os.FileInfo
	path string
}

// NewFileInfoPath creates new FileInfoPath struct
func NewFileInfoPath(f Fatalfable, path string) *FileInfoPath {
	fi, err := os.Stat(path)
	if err != nil {
		f.Fatalf("While getting file %q info: %q", path, err)
	}

	return &FileInfoPath{fi, path}
}

// MakeFipSlice creates a slice of *FileInfoPaths based on
// provided list of file names. It fails on the first
// encountered error.
func MakeFipSlice(f Fatalfable, files ...string) []*FileInfoPath {

	fips := make([]*FileInfoPath, len(files))

	for i, name := range files {
		fips[i] = NewFileInfoPath(f, name)
	}

	return fips
}

// Path returns the stored full path in the FileInfoPath struct
func (fip *FileInfoPath) Path() string {
	return fip.path
}
