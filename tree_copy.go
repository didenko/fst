// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type dirEntry struct {
	name string
	perm os.FileMode
	time time.Time
}

func copyTree(src, dst string) error {

	srcClean := filepath.Clean(src)
	srcLen := len(srcClean)
	dirs := make([]*dirEntry, 0)

	err := filepath.Walk(
		srcClean,
		func(fn string, fi os.FileInfo, er error) error {

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

				_, err = io.Copy(dstf, srcf)
				if err != nil {
					return err
				}

				err = srcf.Close()
				if err != nil {
					return err
				}

				err = dstf.Close()
				if err != nil {
					return err
				}

				err = os.Chmod(dest, fi.Mode())
				if err != nil {
					return err
				}

				destMT := fi.ModTime()
				return os.Chtimes(dest, destMT, destMT)
			}

			if fi.Mode().IsDir() {

				dirs = append(dirs, &dirEntry{dest, fi.Mode().Perm(), fi.ModTime()})
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

		err = os.Chtimes(dirs[i].name, dirs[i].time, dirs[i].time)
		if err != nil {
			return err
		}
	}

	return nil
}
