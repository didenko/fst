// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"
import (
	"os"
	"path/filepath"
	"testing"
)

func TestTreeDiff(t *testing.T) {

	_, cleanup, err := CloneTempChdir("tree_diff_mocks")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	err = filepath.Walk(".", func(p string, i os.FileInfo, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		if filepath.Base(p) == "delete.me" {
			if os.Remove(p) != nil {
				t.Fatal(err)
			}
		}
		return nil
	})

	successes := []string{
		"a_same",
		"d_same_empty",
		"e_same_empty_subdir",
	}

	fails := []string{
		"b_left_nodir", "b_right_nodir",
		"c_left_nofile", "c_right_nofile",
		"f_dir_left_file_right", "f_dir_right_file_left",
		"g_empty_left", "g_empty_right",
		"h_diff_content_bin", "i_diff_content_text_eol",
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
