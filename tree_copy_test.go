// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"log"
	"strings"
	"testing"
)

func TestTreeCopy(t *testing.T) {

	tree := `
		2001-01-01T01:01:01Z	0700	src/
		2001-01-01T01:01:01Z	0550	src/a/
		2099-01-01T01:01:01Z	0700	src/a/b/

		2001-01-01T01:01:01Z	0640	src/c.txt	"This is a two line\nfile with\ta tab\n"
		2001-01-01T01:01:01Z	0600	src/d.txt	No need to quote a single line without tabs

		2002-01-01T01:01:01Z	0700	"src/has\ttab/"
		2002-01-01T01:01:01Z	0440	"src/has\ttab/e.mb"	"# Markdown...\n\n... also ***possible***\n"

		2002-01-01T01:01:01Z	0700	"src/\u10077heavy quoted\u10078/"

		2001-01-01T01:01:01Z	0700	dst/
	`
	nodes, err := TreeParseReader(strings.NewReader(tree))
	if err != nil {
		log.Fatal(err)
	}

	_, cleanup, err := TempCreateChdir(nodes)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	err = TreeCopy("src", "dst")
	if err != nil {
		t.Fatal(err)
	}

	diffs, err := TreeDiff("src", "dst", ByName, ByDir, BySize, ByPerm, ByTime, ByContent(t))
	if err != nil {
		t.Fatal(err)
	}

	if diffs != nil {
		t.Errorf("Trees at \"%s\" and \"%s\" differ unexpectedly: %v", "src", "dst", diffs)
	}
}
