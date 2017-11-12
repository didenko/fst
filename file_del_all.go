// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"os"
	"path/filepath"
)

// FileDelAll recursevely removes file `name` from the `root`
// directory. It is useful to get truly empty directories after
// cloning checked out almost empty directories containing
// a stake file like `.gitkeep`
func FileDelAll(root, name string) error {
	return filepath.Walk(root, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return (err)
		}
		if filepath.Base(p) == name {
			if os.Remove(p) != nil {
				return (err)
			}
		}
		return nil
	})
}
