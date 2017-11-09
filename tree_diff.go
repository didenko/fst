// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"

import (
	"fmt"
	"os"
	"path/filepath"
)

// TreeDiff produces a slice of human-readable notes about
// recursive differences between two directory trees on a
// filesystem. Only plan directories and plain files are
// compared in the tree. Specific comparisons are determined
// By the variadic slice of FileRank functions, like the
// ones in this package. A commonly used set of comparators
// is ByName, ByDir, BySize.
// TODO: , and ByContent
func TreeDiff(a string, b string, comps ...FileRank) []string {
	var diags []string

	listA, err := collectFileInfo(a)
	if err != nil {
		return []string{fmt.Sprintf("Failed to collect entries from \"%s\" with error %v\n", a, err)}
	}

	listB, err := collectFileInfo(b)
	if err != nil {
		return []string{fmt.Sprintf("Failed to collect entries from \"%s\" with error %v\n", b, err)}
	}

	onlyA, onlyB := collectDifferent(listA, listB, comps...)

	if len(onlyA) > 0 {
		diagA := fmt.Sprintf("Unique items from \"%s\": \n", a)
		for _, fi := range onlyA {
			diagA = diagA + fmt.Sprintf("dir:%v, sz:%v, mode:%v, time:%v, name: %v\n", fi.IsDir(), fi.Size(), fi.Mode(), fi.ModTime(), fi.Name())
		}
		diags = append(diags, diagA)
	}
	if len(onlyB) > 0 {
		diagB := fmt.Sprintf("Unique items from \"%s\": \n", b)
		for _, fi := range onlyB {
			diagB = diagB + fmt.Sprintf("dir:%v, sz:%v, mode:%v, time:%v, name: %v\n", fi.IsDir(), fi.Size(), fi.Mode(), fi.ModTime(), fi.Name())
		}
		diags = append(diags, diagB)
	}

	return diags
}

func collectDifferent(left, right []*FileInfoPath, comps ...FileRank) (onlyLeft, onlyRight []*FileInfoPath) {

	onlyLeft = make([]*FileInfoPath, 0)
	onlyRight = make([]*FileInfoPath, 0)

	for l, r := 0, 0; l < len(left) || r < len(right); {

		if r < len(right) && (l == len(left) || Less(right[r], left[l], comps...)) {
			onlyRight = append(onlyRight, right[r])
			r++
			continue
		}

		if l < len(left) && (r == len(right) || Less(left[l], right[r], comps...)) {
			onlyLeft = append(onlyLeft, left[l])
			l++
			continue
		}

		// FIXME: Filenames same, compare content

		l++
		r++
	}

	return onlyLeft, onlyRight
}

func collectFileInfo(dir string) ([]*FileInfoPath, error) {

	list := make([]*FileInfoPath, 0)

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err == nil && path != dir {
			list = append(list, &FileInfoPath{f, path})
		}
		return err
	})

	return list, err
}
