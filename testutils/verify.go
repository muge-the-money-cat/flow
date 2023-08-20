package testutils

import (
	"fmt"

	"github.com/stretchr/testify/assert"
)

func Verify(assertion assertionFunc, expected, actual interface{}) (e error) {
	var (
		t testingT
	)

	assertion(&t, expected, actual)

	return t.e
}

type assertionFunc func(
	t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{},
) bool

type testingT struct {
	e error
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	t.e = fmt.Errorf(format, args...)
}
