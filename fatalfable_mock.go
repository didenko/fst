// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"fmt"
	"strings"
	"testing"
)

type testFatalfCatch struct {
	*testing.T
	caught bool
	prefix string
}

func wrapFatalfTest(t *testing.T, messageStart string) *testFatalfCatch {
	wrapper := &testFatalfCatch{t, false, messageStart}
	wrapper.Cleanup(wrapper.failIfPassed)
	return wrapper
}

func (tef *testFatalfCatch) Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if strings.HasPrefix(message, tef.prefix) {
		tef.Logf("Fatalf called as expected: "+format, v...)
		tef.caught = true
		tef.SkipNow()
	} else {
		tef.T.Fatalf(format, v...)
	}
}

func (tef *testFatalfCatch) failIfPassed() {
	if !tef.caught {
		tef.T.Fatalf("Test '%s' did not call Fatalf('%s'...) as expected in the tested function", tef.Name(), tef.prefix)
	}
}
