// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"os"
	"time"
)

// Node holds basic attributes of a filesystem item.
// Its name is relative to CWD.
type Node struct {
	perm os.FileMode
	time time.Time
	name string
	body string
}

// SaveAttributes sets the named file's permissions and
// timestamps to the ones from the node.
func (n *Node) SaveAttributes(f Fatalfable) {

	err := os.Chmod(n.name, n.perm)
	if err != nil {
		f.Fatalf("Setting %q permissions to %o: %q", n.name, n.perm, err)
	}

	err = os.Chtimes(n.name, n.time, n.time)
	if err != nil {
		f.Fatalf("Setting %q timestamps to %s: %q", n.name, n.time, err)
	}
}

// Rfc3339 converts a string to a time struct while assuming
// the string is formatted according to RFC3339. It calls
// f.Fatalf if the conversion fails.
func Rfc3339(f Fatalfable, ts string) time.Time {
	tm, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		f.Fatalf("Failed to convert %q to a time: %q", ts, err)
	}
	return tm
}
