// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// TreeCopy duplicates regular files and directories from
// inside the source directory into an existing destination
// directory.
func TreeCopy(f Fatalfable, src, dst string) {

	srcClean := filepath.Clean(src)
	srcLen := len(srcClean)
	dirs := make([]*Node, 0)

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
					return fmt.Errorf("Opening the sorce file %q: %s", fn, err)
				}

				dstf, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0600)
				if err != nil {
					return fmt.Errorf("Opening the dest file %q: %s", dest, err)
				}

				_, err = io.Copy(dstf, srcf)
				if err != nil {
					return fmt.Errorf("Copying %q to %q: %s", fn, dest, err)
				}

				err = srcf.Close()
				if err != nil {
					return fmt.Errorf("Closing the sorce file %q: %s", fn, err)
				}

				err = dstf.Close()
				if err != nil {
					return fmt.Errorf("Closing the dest file %q: %s", dest, err)
				}

				err = os.Chmod(dest, fi.Mode())
				if err != nil {
					return fmt.Errorf("Setting permissions on %q: %s", dest, err)
				}

				destMT := fi.ModTime()
				err = os.Chtimes(dest, destMT, destMT)
				if err != nil {
					return fmt.Errorf("Setting timestamp %s on %q: %s", destMT, dest, err)
				}

				return nil
			}

			if fi.Mode().IsDir() {

				dirs = append(dirs, &Node{fi.Mode().Perm(), fi.ModTime(), dest, ""})
				err := os.Mkdir(dest, 0700)
				if err != nil {
					return fmt.Errorf("Creating dir %q: %s", dest, err)
				}
			}
			return nil
		})

	if err != nil {
		f.Fatalf("Copying tree from %q to %q: %s", src, dst, err)
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		err := os.Chmod(dirs[i].name, dirs[i].perm)
		if err != nil {
			f.Fatalf("Setting permissions on %q to %o: %s", dirs[i].name, dirs[i].perm, err)
		}

		err = os.Chtimes(dirs[i].name, dirs[i].time, dirs[i].time)
		if err != nil {
			f.Fatalf("Setting timestamp %s on %q: %s", dirs[i].time, dirs[i].name, err)
		}
	}
}
