// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"
import "testing"

func TestTreeDiff(t *testing.T) {

	_, cleanup, err := CloneTempChdir("tree_diff_mocks")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	if diffs := TreeDiff("a_same/left", "a_same/right", ByName, ByDir, BySize, ByContent(t)); diffs != nil {
		t.Errorf("Equivalent directories in \"a_same\" tested as different: %v\n", diffs)
	}

	if diffs := TreeDiff("ba_left_nodir/left", "ba_left_nodir/right", ByName, ByDir, BySize, ByContent(t)); diffs == nil {
		t.Error("Differing directories in \"ba_left_nodir\" tested as equivalent\n")
	}
}
