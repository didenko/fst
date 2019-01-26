// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"testing"
)

func TestFileDelAll(t *testing.T) {

	nodes := []*Node{
		&Node{0750, Rfc3339(t, "2018-01-01T01:01:01Z"), "mock/", ""},
		&Node{0640, Rfc3339(t, "2018-01-02T03:04:05Z"), "mock/.gitkeep", ""},
		&Node{0640, Rfc3339(t, "2018-02-03T04:05:06Z"), "mock/text.txt", ""},
		&Node{0750, Rfc3339(t, "2018-03-04T05:06:07Z"), "mock/dir/", ""},
		&Node{0640, Rfc3339(t, "2018-04-05T06:07:08Z"), "mock/dir/.gitkeep", ""},
		&Node{0640, Rfc3339(t, "2018-05-06T07:08:09Z"), "mock/dir/text.txt", ""},
		&Node{0750, Rfc3339(t, "2018-06-04T08:09:10Z"), "mock/dir/dir/", ""},
		&Node{0640, Rfc3339(t, "2018-07-05T09:10:11Z"), "mock/dir/dir/.gitkeep", ""},
		&Node{0640, Rfc3339(t, "2018-08-06T10:11:12Z"), "mock/dir/dir/text.txt", ""},
		&Node{0750, Rfc3339(t, "2018-01-01T01:01:01Z"), "test/", ""},
		&Node{0640, Rfc3339(t, "2018-02-03T04:05:06Z"), "test/text.txt", ""},
		&Node{0750, Rfc3339(t, "2018-03-04T05:06:07Z"), "test/dir/", ""},
		&Node{0640, Rfc3339(t, "2018-05-06T07:08:09Z"), "test/dir/text.txt", ""},
		&Node{0750, Rfc3339(t, "2018-06-04T08:09:10Z"), "test/dir/dir/", ""},
		&Node{0640, Rfc3339(t, "2018-08-06T10:11:12Z"), "test/dir/dir/text.txt", ""},
	}

	_, cleanup := TempCreateChdir(t, nodes)
	defer cleanup()

	FileDelAll(t, "mock", ".gitkeep")

	diffs, err := TreeDiff(t, "mock", "test", ByName, ByDir, BySize)
	if err != nil {
		t.Fatal(err)
	}

	if len(diffs) > 0 {
		t.Errorf("Tree foes not match expected after FileDelAll:\n%v\n", diffs)
	}
}
