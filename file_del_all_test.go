// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import "testing"
import "strings"

func TestFileDelAll(t *testing.T) {

	tree := `
		2018-01-01T01:01:01Z	0750	mock/
		2018-01-02T03:04:05Z	0640	mock/.gitkeep
		2018-02-03T04:05:06Z	0640	mock/text.txt
		2018-03-04T05:06:07Z	0750	mock/dir/
		2018-04-05T06:07:08Z	0640	mock/dir/.gitkeep
		2018-05-06T07:08:09Z	0640	mock/dir/text.txt
		2018-06-04T08:09:10Z	0750	mock/dir/dir/
		2018-07-05T09:10:11Z	0640	mock/dir/dir/.gitkeep
		2018-08-06T10:11:12Z	0640	mock/dir/dir/text.txt
		2018-01-01T01:01:01Z	0750	test/
		2018-02-03T04:05:06Z	0640	test/text.txt
		2018-03-04T05:06:07Z	0750	test/dir/
		2018-05-06T07:08:09Z	0640	test/dir/text.txt
		2018-06-04T08:09:10Z	0750	test/dir/dir/
		2018-08-06T10:11:12Z	0640	test/dir/dir/text.txt
	`
	nodes, err := TreeParseReader(strings.NewReader(tree))
	if err != nil {
		t.Fatal(err)
	}

	_, cleanup, err := TempCreateChdir(nodes)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	err = FileDelAll("mock", ".gitkeep")
	if err != nil {
		t.Fatal(err)
	}

	diffs, err := TreeDiff("mock", "test", ByName, ByDir, BySize)
	if err != nil {
		t.Fatal(err)
	}

	if len(diffs) > 0 {
		t.Errorf("Tree foes not match expected after FileDelAll:\n%v\n", diffs)
	}
}
