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
package fstests // import "go.didenko.com/fstests"

import (
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
