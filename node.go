// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"os"
	"time"
)

// Node holds basic attributes of a filesystem item.
// Its name is relative to CWD.
type Node struct {
	name string
	perm os.FileMode
	time time.Time
	body string
}
