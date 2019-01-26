// Copyright 2017-2019 Vlad Didenko. All rights reserved.
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
func TreeCreate(f Fatalfable, entries []*Node) {
	dirs := make([]*Node, 0)

	for _, e := range entries {

		if e.name[len(e.name)-1] == '/' {
			if err := os.Mkdir(e.name[:len(e.name)-1], 0700); err != nil {
				f.Fatalf("While making dir %q: %s", e.name, err)
			}

			dirs = append(dirs, e)
			continue
		}

		fl, err := os.Create(e.name)
		if err != nil {
			f.Fatalf("While creating the file %q: %s", e.name, err)
		}

		if len(e.body) > 0 {
			_, err = fl.WriteString(e.body)
			if err != nil {
				f.Fatalf("While writing file %q content: %s", e.name, err)
			}
		}

		err = fl.Close()
		if err != nil {
			f.Fatalf("While colsing file %q: %s", e.name, err)
		}

		e.SaveAttributes(f)
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		dirs[i].SaveAttributes(f)
	}
}
