// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"os"
	"testing"
)

func TestFileDelAll(t *testing.T) {

	nodes := []*Node{
		{0750, Rfc3339(t, "2018-01-01T01:01:01Z"), "mock/", ""},
		{0640, Rfc3339(t, "2018-01-02T03:04:05Z"), "mock/.gitkeep", ""},
		{0640, Rfc3339(t, "2018-02-03T04:05:06Z"), "mock/text.txt", ""},
		{0750, Rfc3339(t, "2018-03-04T05:06:07Z"), "mock/dir/", ""},
		{0640, Rfc3339(t, "2018-04-05T06:07:08Z"), "mock/dir/.gitkeep", ""},
		{0640, Rfc3339(t, "2018-05-06T07:08:09Z"), "mock/dir/text.txt", ""},
		{0750, Rfc3339(t, "2018-06-04T08:09:10Z"), "mock/dir/dir/", ""},
		{0640, Rfc3339(t, "2018-07-05T09:10:11Z"), "mock/dir/dir/.gitkeep", ""},
		{0640, Rfc3339(t, "2018-08-06T10:11:12Z"), "mock/dir/dir/text.txt", ""},
		{0750, Rfc3339(t, "2018-01-01T01:01:01Z"), "test/", ""},
		{0640, Rfc3339(t, "2018-02-03T04:05:06Z"), "test/text.txt", ""},
		{0750, Rfc3339(t, "2018-03-04T05:06:07Z"), "test/dir/", ""},
		{0640, Rfc3339(t, "2018-05-06T07:08:09Z"), "test/dir/text.txt", ""},
		{0750, Rfc3339(t, "2018-06-04T08:09:10Z"), "test/dir/dir/", ""},
		{0640, Rfc3339(t, "2018-08-06T10:11:12Z"), "test/dir/dir/text.txt", ""},
	}

	_, cleanup := TempCreateChdir(t, nodes)
	defer cleanup()

	FileDelAll(t, "mock", ".gitkeep")

	diffs := TreeDiff(t, "mock", "test", ByName, ByDir, BySize)

	if len(diffs) > 0 {
		t.Errorf("Tree foes not match expected after FileDelAll:\n%v\n", diffs)
	}
}

func TestFileDelAllWalkError(t *testing.T) {
	nodes := []*Node{
		{0750, Rfc3339(t, "2018-01-01T01:01:01Z"), "mock/", ""},
		{0750, Rfc3339(t, "2018-03-04T05:06:07Z"), "mock/dir/", ""},
		{0640, Rfc3339(t, "2018-04-05T06:07:08Z"), "mock/dir/.gitkeep", ""},
		{0640, Rfc3339(t, "2018-02-03T04:05:06Z"), "mock/dir/text.txt", ""},
	}

	_, cleanup := TempCreateChdir(t, nodes)
	defer cleanup()

	if err := os.Chmod("mock/dir", 0644); err != nil {
		t.Errorf("Test setup failed to change directory permissions to u-x: %v", err)
	}

	testWrapper := wrapFatalfTest(t, "Remove: while walking to \"mock/dir/.gitkeep\": ")
	FileDelAll(testWrapper, "mock", ".gitkeep")
}

func TestFileDelAllRemoveError(t *testing.T) {
	nodes := []*Node{
		{0750, Rfc3339(t, "2018-01-01T01:01:01Z"), "mock/", ""},
		{0750, Rfc3339(t, "2018-03-04T05:06:07Z"), "mock/dir/", ""},
		{0640, Rfc3339(t, "2018-04-05T06:07:08Z"), "mock/dir/.gitkeep", ""},
		{0640, Rfc3339(t, "2018-02-03T04:05:06Z"), "mock/dir/text.txt", ""},
	}

	_, cleanup := TempCreateChdir(t, nodes)
	defer cleanup()

	if err := os.Chmod("mock/dir", 0550); err != nil {
		t.Errorf("Test setup failed to change directory permissions to u-w: %v", err)
	}

	testROWrap := wrapFatalfTest(t, "Removing \"mock/dir/.gitkeep\": ")
	FileDelAll(testROWrap, "mock", ".gitkeep")
}
