// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// TempInitDir function creates a directory for holding
// temporary files according to platform preferences and
// returns the directory name and a cleanup function.
//
// The returned values are:
//
// 1. a string containing the created temporary directory path
//
// 2. a cleanup function to change back to the old working
// directory and to delete the temporary directory
//
// 3. an error
//
// If there was an error while creating the temporary
// directory, then the returned directory name is empty,
// cleanup funcion is nil, and the temp folder is
// expected to be already removed.
func TempInitDir(f Fatalfable) (string, func()) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		os.RemoveAll(root)
		f.Fatalf("Creating a temp directory: %s", err)
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
			f.Fatalf("In the file cleanup: %s", err)
		}

		for i := len(dirs) - 1; i >= 0; i-- {
			err = os.RemoveAll(dirs[i])
			if err != nil {
				f.Fatalf("In the directory cleanup: %s", err)
			}
		}
	}
}

// TempInitChdir creates a temporary directory in the same
// fashion as TempInitDir. It also changes into the newly
// created temporary directory and adds returning back
// to the old working directory to the returned cleanup
// function. The returned values are:
//
// 1. a string containing the previous working directory
//
// 2. a cleanup function to change back to the old working
// directory and to delete the temporary directory
//
// 3. an error
func TempInitChdir(f Fatalfable) (string, func()) {
	root, cleanup := TempInitDir(f)

	wd, err := os.Getwd()
	if err != nil {
		cleanup()
		f.Fatalf("Working directory indetermined: %s", err)
	}

	err = os.Chdir(root)
	if err != nil {
		cleanup()
		f.Fatalf("Changing to the temporary directory %q: %s", root, err)
	}

	return wd,
		func() {
			os.Chdir(wd)
			cleanup()
		}
}

// TempCloneDir function creates a copy of an existing
// directory with it's content - regular files only - for
// holding temporary test files.
//
// The returned values are:
//
// 1. a string containing the created temporary directory path
//
// 2. a cleanup function to change back to the old working
// directory and to delete the temporary directory
//
// 3. an error
//
// If there was an error while cloning the temporary
// directory, then the returned directory name is empty,
// cleanup funcion is nil, and the temp folder is
// expected to be already removed.
//
// The clone attempts to maintain the basic original Unix
// permissions (9-bit only, from the rxwrwxrwx set).
// If, however, the user does not have read permission
// for a file, or read+execute permission for a directory,
// then the clone process will naturally fail.
func TempCloneDir(f Fatalfable, src string) (string, func()) {
	root, cleanup := TempInitDir(f)
	TreeCopy(newFatalCleaner(f, cleanup), src, root)
	return root, cleanup
}

// TempCloneChdir clones a temporary directory in the same
// fashion as TempCloneDir. It also changes into the newly
// cloned temporary directory and adds returning back
// to the old working directory to the returned cleanup
// function. The returned values are:
//
// 1. a string containing the previous working directory
//
// 2. a cleanup function to change back to the old working
// directory and to delete the temporary directory
//
// 3. an error
func TempCloneChdir(f Fatalfable, src string) (string, func(), error) {
	root, cleanup := TempCloneDir(f, src)

	wd, err := os.Getwd()
	if err != nil {
		cleanup()
		f.Fatalf("Working directory indetermined: %s", err)
	}

	err = os.Chdir(root)
	if err != nil {
		cleanup()
		f.Fatalf("Changing to the temporary directory %q: %s", root, err)
	}

	return wd,
		func() {
			os.Chdir(wd)
			cleanup()
		},
		nil
}

// TempCreateChdir is a combination of `TempInitChdir` and
// `TreeCreate` functions. It creates a termporary directory,
// changes into it, populates it fron the provided `config`
// as `TreeCreate` would, and returns the old directory name
// and the cleanup function.
func TempCreateChdir(f Fatalfable, nodes []*Node) (string, func()) {
	old, cleanup := TempInitChdir(f)
	TreeCreate(newFatalCleaner(f, cleanup), nodes)
	return old, cleanup
}
