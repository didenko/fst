// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"
import (
	"os"
)

// TreeCreate creates the filesystem objects provided in the
// slice of Node pointers, where Nodes describe the objects
// to be created.
//
// It is up to the caller to deal with conflicting file and
// directory names in the input. TreeCreate processes
// the input line-by-line and will return with error at a first
// problem it runs into.
func TreeCreate(entries []*Node) error {

	dirs := make([]*Node, 0)

	for _, e := range entries {

		if e.name[len(e.name)-1] == '/' {
			if err := os.Mkdir(e.name[:len(e.name)-1], 0700); err != nil {
				return err
			}

			dirs = append(dirs, e)
			continue
		}

		f, err := os.Create(e.name)
		if err != nil {
			return err
		}

		if len(e.body) > 0 {
			_, err = f.WriteString(e.body)
			if err != nil {
				return err
			}
		}

		err = f.Close()
		if err != nil {
			return err
		}

		if err = e.SaveAttributes(); err != nil {
			return err
		}
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		if err := dirs[i].SaveAttributes(); err != nil {
			return err
		}
	}

	return nil
}
