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

type dirEntry struct {
	name string
	perm os.FileMode
}

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
		dirs := make([]string, 0)

		err := filepath.Walk(
			root,
			func(fn string, fi os.FileInfo, er error) error {

				if fi.IsDir() {
					err = os.Chmod(fn, 0700)
					if err != nil {
						return err
					}

					dirs = append(dirs, fn)
					return nil
				}

				return os.Remove(fn)
			})

		if err != nil {
			log.Fatalln(err)
		}

		for i := len(dirs) - 1; i >= 0; i-- {
			err = os.RemoveAll(dirs[i])
			if err != nil {
				log.Fatalln(err)
			}
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
//
// The clone attempts to maintain the basic original Unix
// permissions (9-bit only, from the rxwrwxrwx set).
// If, however, the user does not have read permission
// for a file, or read+execute permission for a directory,
// then the clone process will naturally fail.
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
	dirs := make([]*dirEntry, 0)

	err := filepath.Walk(
		srcClean,
		func(fn string, fi os.FileInfo, er error) error {

			// TODO: set proper permissions, including 10-12 bits
			// TODO: set proper ownership
			// TODO: set proper timestamps

			if er != nil || len(fn) <= srcLen {
				return er
			}

			dest := filepath.Join(dst, fn[srcLen:])

			if fi.Mode().IsRegular() {

				srcf, err := os.Open(fn)
				if err != nil {
					return err
				}

				dstf, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0600)
				if err != nil {
					return err
				}

				err = dstf.Chmod(fi.Mode().Perm())
				if err != nil {
					return err
				}

				_, err = io.Copy(dstf, srcf)
				return err
			}

			if fi.IsDir() {
				dirs = append(dirs, &dirEntry{dest, fi.Mode().Perm()})
				return os.Mkdir(dest, 0700)
			}

			return nil
		})

	if err != nil {
		return err
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		err := os.Chmod(dirs[i].name, dirs[i].perm)
		if err != nil {
			return err
		}
	}

	return nil
}

// TreeDiffs produces a slice of human-readable notes about
// recursive differences between two directory trees on a
// filesystem. Only plan directories and plain files are
// compared in the tree. The follwing attributes are compared:
//
// 1. Name
// 2. Size
// 3. Permissions (only Unix's 9 bits)
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

		// FIXME: Filenames same, compare content

		l++
		r++
	}

	return onlyLeft, onlyRight
}

func less(left, right os.FileInfo) bool {

	// TODO: test high (10-12) permission bits
	// TODO: test ownership
	// TODO: test timestamps

	return left.Name() < right.Name() ||
		left.IsDir() != right.IsDir() ||
		left.Size() < right.Size() ||
		left.Mode().Perm() < right.Mode().Perm()
}
