// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"testing"
)

func TestTreeCopy(t *testing.T) {

	nodes := []*Node{
		&Node{0700, Rfc3339(t, "2001-01-01T01:01:01Z"), "src/", ""},
		&Node{0550, Rfc3339(t, "2001-01-01T01:01:01Z"), "src/a/", ""},
		&Node{0700, Rfc3339(t, "2099-01-01T01:01:01Z"), "src/a/b/", ""},
		&Node{0640, Rfc3339(t, "2001-01-01T01:01:01Z"), "src/c.txt", "This is a two line\nfile with\ta tab\n"},
		&Node{0600, Rfc3339(t, "2001-01-01T01:01:01Z"), "src/d.txt", "A single line without tabs"},
		&Node{0700, Rfc3339(t, "2002-01-01T01:01:01Z"), "src/has\ttab/", ""},
		&Node{0440, Rfc3339(t, "2002-01-01T01:01:01Z"), "src/has\ttab/e.mb", "# Markdown...\n\n... also ***possible***\n"},
		&Node{0700, Rfc3339(t, "2002-01-01T01:01:01Z"), "src/\u10077heavy quoted\u10078/", ""},
		&Node{0700, Rfc3339(t, "2001-01-01T01:01:01Z"), "dst/", ""},
	}

	_, cleanup := TempCreateChdir(t, nodes)
	defer cleanup()

	TreeCopy(t, "src", "dst")

	diffs := TreeDiff(t, "src", "dst", ByName, ByDir, BySize, ByPerm, ByTime, ByContent(t))

	if diffs != nil {
		t.Errorf("Trees at \"%s\" and \"%s\" differ unexpectedly: %v", "src", "dst", diffs)
	}
}
