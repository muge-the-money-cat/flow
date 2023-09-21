package testutils

import (
	"fmt"

	"github.com/stretchr/testify/assert"
)

func Verify(assertion assertionFunc, iface0, iface1 interface{}) (e error) {
	var (
		t testingT
	)

	assertion(&t, iface0, iface1)

	return t.e
}

type assertionFunc func(
	assert.TestingT, interface{}, interface{}, ...interface{},
) bool

type testingT struct {
	e error
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	t.e = fmt.Errorf(format, args...)
}
