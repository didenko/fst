// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"
import (
	"path/filepath"
	"testing"
)

func TestTreeDiff(t *testing.T) {

	_, cleanup, err := CloneTempChdir("tree_diff_mocks")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	successes := []string{"a_same"}
	fails := []string{
		"b_left_nodir", "b_right_nodir",
		"c_left_nofile", "c_right_nofile",
	}

	for _, caseDir := range successes {
		if diffs := TreeDiff(
			filepath.Join(caseDir, "left"),
			filepath.Join(caseDir, "right"),
			ByName, ByDir, BySize, ByContent(t),
		); diffs != nil {
			t.Errorf("Equivalent directories in \"%s\" tested as different: %v\n", caseDir, diffs)
		}
	}

	for _, caseDir := range fails {
		if diffs := TreeDiff(
			filepath.Join(caseDir, "left"),
			filepath.Join(caseDir, "right"),
			ByName, ByDir, BySize, ByContent(t),
		); diffs == nil {
			t.Errorf("Differing directories in \"%s\" passed as equivalent\n", caseDir)
		}
	}
}
