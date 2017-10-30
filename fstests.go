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
package fstests // import "didenko.com/go/fstests"

import (
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

			switch mode := fi.Mode(); mode {

			case os.ModeDir:
				// FIXME: set a proper mode
				return os.MkdirAll(dest, 0750)

			case os.ModeDevice,
				os.ModeNamedPipe,
				os.ModeSocket,
				os.ModeSymlink,
				os.ModeTemporary:
				return nil

			default:
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
		})
}
