// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"os"
	"time"
)

// DirEntry holds basic attributes of a filesystem item.
// Its name is relative to CWD.
type DirEntry struct {
	name string
	perm os.FileMode
	time time.Time
	body string
}
