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

	type DiffCase struct {
		dir   string
		comps []FileRank
	}

	successes := []DiffCase{
		DiffCase{"a_same_content", []FileRank{ByName, ByDir, BySize, ByContent(t)}},
		DiffCase{"d_same_empty", []FileRank{ByName, BySize}},
		DiffCase{"e_same_empty_subdir", []FileRank{ByName, BySize}},
		DiffCase{"k_same_size", []FileRank{ByName, BySize}},
		DiffCase{"j_diff_sizes_same_perm", []FileRank{ByName, ByPerm}},
		DiffCase{"l_perms_same", []FileRank{ByName, ByPerm}},
	}

	fails := []DiffCase{
		DiffCase{"b_left_nodir", []FileRank{ByName}},
		DiffCase{"b_right_nodir", []FileRank{ByName}},
		DiffCase{"c_left_nofile", []FileRank{ByName}},
		DiffCase{"c_right_nofile", []FileRank{ByName}},
		DiffCase{"f_dir_left_file_right", []FileRank{ByName, ByDir}},
		DiffCase{"f_dir_right_file_left", []FileRank{ByName, ByDir}},
		DiffCase{"g_empty_left", []FileRank{ByName}},
		DiffCase{"g_empty_right", []FileRank{ByName}},
		DiffCase{"h_diff_content_bin", []FileRank{ByName, ByContent(t)}},
		DiffCase{"i_diff_content_text_eol", []FileRank{ByName, ByContent(t)}},
		DiffCase{"j_diff_sizes_same_perm", []FileRank{ByName, BySize}},
		DiffCase{"l_perms_same", []FileRank{ByName, ByPerm, BySize}},
	}

	for _, tc := range successes {
		if diffs := TreeDiff(
			filepath.Join(tc.dir, "left"),
			filepath.Join(tc.dir, "right"),
			tc.comps...,
		); diffs != nil {
			t.Errorf("Equivalent directories in \"%s\" tested as different: %v\n", tc.dir, diffs)
		}
	}

	for _, tc := range fails {
		if diffs := TreeDiff(
			filepath.Join(tc.dir, "left"),
			filepath.Join(tc.dir, "right"),
			tc.comps...,
		); diffs == nil {
			t.Errorf("Differing directories in \"%s\" passed as equivalent\n", tc.dir)
		}
	}
}
