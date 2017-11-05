// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TreeDiffs produces a slice of human-readable notes about
// recursive differences between two directory trees on a
// filesystem. Only plan directories and plain files are
// compared in the tree. The following attributes are compared:
//
// 1. Name
//
// 2. Size
//
// 3. Permissions (only Unix's 12 bits)
//
// 4. Timestamps
func TreeDiffs(a string, b string) []string {
	var diags []string

	listA, err := collectFileInfo(a)
	if err != nil {
		return []string{fmt.Sprintf("Failed to collect entries from \"%s\" with error %v\n", a, err)}
	}

	listB, err := collectFileInfo(b)
	if err != nil {
		return []string{fmt.Sprintf("Failed to collect entries from \"%s\" with error %v\n", b, err)}
	}

	onlyA, onlyB := collectDifferent(listA, listB)

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

func collectDifferent(left, right []os.FileInfo) (onlyLeft, onlyRight []os.FileInfo) {

	onlyLeft = make([]os.FileInfo, 0)
	onlyRight = make([]os.FileInfo, 0)

	for l, r := 0, 0; l < len(left) || r < len(right); {

		if r < len(right) && (l == len(left) || less(right[r], left[l])) {
			onlyRight = append(onlyRight, right[r])
			r++
			continue
		}

		if l < len(left) && (r == len(right) || less(left[l], right[r])) {
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

func less(left, right os.FileInfo) bool {

	return left.Name() < right.Name() ||
		(left.IsDir() && !right.IsDir()) ||
		(!left.IsDir() && (left.Size() < right.Size())) ||
		left.Mode() < right.Mode() ||
		left.ModTime().Before(right.ModTime().Add(-5*time.Millisecond))
}

func collectFileInfo(dir string) ([]os.FileInfo, error) {

	list := []os.FileInfo{}

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err == nil && path != dir {
			list = append(list, f)
		}
		return err
	})

	return list, err
}
