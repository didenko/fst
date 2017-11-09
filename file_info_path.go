// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"

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
func NewFileInfoPath(path string) (*FileInfoPath, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return &FileInfoPath{nil, path}, err
	}

	return &FileInfoPath{fi, path}, nil
}

// MakeFipSlice creates a slice of *FileInfoPaths based on
// provided list of file names. It returns the first
// encountered error.
func MakeFipSlice(files ...string) ([]*FileInfoPath, error) {

	fips := make([]*FileInfoPath, 0)

	for _, name := range files {

		fip, err := NewFileInfoPath(name)
		if err != nil {
			return nil, err
		}

		fips = append(fips, fip)
	}

	return fips, nil
}

// Path returns the stored full path in the FileInfoPath struct
func (fip *FileInfoPath) Path() string {
	return fip.path
}
