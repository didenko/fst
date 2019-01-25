// Copyright 2017-2019 Vlad Didenko. All rights reserved.
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
func FileDelAll(f Fatalfable, root, name string) {
	err := filepath.Walk(
		root,
		func(p string, i os.FileInfo, err error) error {
			if err != nil {
				f.Fatalf("Remove: while walking to %q: %q", p, err)
			}
			if filepath.Base(p) == name {
				if os.Remove(p) != nil {
					f.Fatalf("Removing %q: %q", p, err)
				}
			}
			return nil
		},
	)
	if err != nil {
		f.Fatalf("Removing %q: %q", root, err)
	}
}
