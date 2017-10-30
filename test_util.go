// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstests // import "didenko.com/go/fstests"

import (
	"log"
	"os"
	"path/filepath"
)

func treesAreDifferent(a string, b string) bool {

	listA, err := collectFileInfo(a)
	if err != nil {
		log.Printf("Failed to collect entries from \"%s\" with error %v\n", a, err)
		return true
	}

	listB, err := collectFileInfo(b)
	if err != nil {
		log.Printf("Failed to collect entries from \"%s\" with error %v\n", b, err)
		return true
	}

	ret := false

	onlyA, onlyB := collectDifferent(listA, listB)

	if len(onlyA) > 0 {
		log.Printf("Unique items from \"%s\": %v\n", a, onlyA)
		ret = true
	}
	if len(onlyB) > 0 {
		log.Printf("Unique items from \"%s\": %v\n", b, onlyB)
		ret = true
	}

	return ret
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

func less(left, right os.FileInfo) bool {
	return left.Name() < right.Name() ||
		left.IsDir() != right.IsDir() ||
		left.Size() < right.Size()
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

		// FIXME: Filenames same, compare:
		// TODO: size
		// TODO: content
		// TODO: permissions?
		// TODO: ACLs?

		l++
		r++
	}

	return onlyLeft, onlyRight
}
