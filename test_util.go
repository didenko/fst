// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstests // import "didenko.com/go/fstests"

import (
	"os"
	"path/filepath"
)

func collectFileInfo(dir string) ([]os.FileInfo, error) {

	list := []os.FileInfo{}

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err == nil && path != dir {
			list = append(list, f)
		}
		return err
	})

	return list, err
}

func less(left, right os.FileInfo) bool {
	return left.Name() < right.Name() ||
		left.IsDir() != right.IsDir() ||
		left.Size() < right.Size()
}
