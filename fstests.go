// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

// Package fstests is a collection of functions to help
// testing filesyste objects modifications. It focuses
// on creating and cleaning up a baseline filesyetem state.
//
// The following are addressed use cases:
//
// 1. Create a directory hierarchy via an API
// 2. Create a directory hierarchy via a copy of a template
// 3. Write a provided test data to files
// 4. Contain all test activity in a temporatry directory
// 5. Compare two directories recursively
package fstests // import "didenko.com/go/fstests"

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var noop = func() {}

// InitTempDir function creates a directory for holding
// temporary files according to platform preferences and
// returns the directory name and a cleanup function.
//
// If there was an error while creating the temporary
// directory, then the returned directory name is empty,
// cleanup funcion is a noop, and the temp folder is
// expected to be already removed.
func InitTempDir() (string, func(), error) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		os.RemoveAll(root)
		return "", noop, err
	}

	return root, func() {
		err := os.RemoveAll(root)
		if err != nil {
			log.Fatalf("Error while removing the temporary directory %s", root)
		}
	}, nil
}

// CloneTempDir function creates a copy of an existing
// directory with it's content - regular files only - for
// holding temporary test files. It returns the directory
// name and a cleanup function.
//
// If there was an error while cloning the temporary
// directory, then the returned directory name is empty,
// cleanup funcion is a noop, and the temp folder is
// expected to be already removed.
func CloneTempDir(src string) (string, func(), error) {
	root, cleanup, err := InitTempDir()
	if err != nil {
		return "", noop, err
	}

	err = copyTree(src, root)
	if err != nil {
		cleanup()
		return "", noop, err
	}

	return root, cleanup, nil
}

func copyTree(src, dst string) error {

	srcClean := filepath.Clean(src)
	srcLen := len(srcClean)

	return filepath.Walk(
		srcClean,
		func(fn string, fi os.FileInfo, er error) error {

			if er != nil || len(fn) <= srcLen {
				return er
			}

			dest := filepath.Join(dst, fn[srcLen:])

			if fi.Mode().IsRegular() {
				// FIXME: set a proper mode
				srcf, err := os.Open(fn)
				if err != nil {
					return err
				}

				dstf, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0640)
				if err != nil {
					return err
				}

				_, err = io.Copy(dstf, srcf)
				return err
			}

			if fi.IsDir() {
				return os.MkdirAll(dest, 0750)
			}

			return nil
		})
}

// TreeDiffs produces a slice of human-readable notes about
// recursive differences between two directory trees on a
// filesystem. Only plan directories and plain files are
// compared in the tree. The follwing attributes are compared:
//
// 1. Name
// 2. Size
// 3. Permissions
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
		diags = append(diags, fmt.Sprintf("Unique items from \"%s\": %v\n", a, onlyA))
	}
	if len(onlyB) > 0 {
		diags = append(diags, fmt.Sprintf("Unique items from \"%s\": %v\n", b, onlyB))
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

		// FIXME: Filenames same, compare:
		// TODO: content
		// TODO: permissions?
		// TODO: ACLs?

		l++
		r++
	}

	return onlyLeft, onlyRight
}
