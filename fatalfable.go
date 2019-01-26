// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

// Fatalfable is an interface to any type containing a common
// Fatalf method, as the likes of testing.T and log.Logger.
type Fatalfable interface {
	Fatalf(format string, v ...interface{})
}

type fatalCleaner struct {
	utter Fatalfable
	clean func()
}

func newFatalCleaner(f Fatalfable, c func()) *fatalCleaner {
	return &fatalCleaner{f, c}
}

func (fc *fatalCleaner) Fatalf(f string, args ...interface{}) {
	fc.clean()
	fc.utter.Fatalf(f, args...)
}
