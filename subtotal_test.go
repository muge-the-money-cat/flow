package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/stretchr/testify/assert"
)

func TestPostSubtotalWithNoParent(t *testing.T) {
	var (
		testSuite = godog.TestSuite{
			ScenarioInitializer: initScenarioPostSubtotalWithNoParent,
			Options: &godog.Options{
				Output: colors.Colored(os.Stdout),
				Format: "pretty",
			},
		}

		status int
	)

	status = testSuite.Run()

	if status != 0 {
		os.Exit(status)
	}

	return
}

func initScenarioPostSubtotalWithNoParent(ctx *godog.ScenarioContext) {
	ctx.Step(`^there is a Subtotal API$`,
		thereIsASubtotalAPI,
	)
	ctx.Step(`^we POST a Subtotal with name "(.+)" and no parent$`,
		postASubtotalWithNoParent,
	)
	ctx.Step(`^there should be a Subtotal with name "(.+)" and no parent$`,
		thereShouldBeASubtotalWithNoParent,
	)

	return
}

func thereIsASubtotalAPI(parentContext context.Context) (
	childContext context.Context, e error,
) {
	var (
		api SubtotalAPI = NewInMemorySubtotalAPI()
	)

	childContext = context.WithValue(parentContext,
		subtotalAPIContextKey{},
		api,
	)

	return
}

func postASubtotalWithNoParent(parentContext context.Context, name string,
) (
	childContext context.Context, e error,
) {
	var (
		subtotal Subtotal = NewSubtotalWithNoParent(name)
	)

	childContext = parentContext

	e = childContext.Value(subtotalAPIContextKey{}).(SubtotalAPI).Post(subtotal)
	if e != nil {
		return
	}

	return
}

func thereShouldBeASubtotalWithNoParent(
	parentContext context.Context, name string,
) (
	childContext context.Context, e error,
) {
	var (
		subtotal Subtotal
	)

	childContext = parentContext

	subtotal, e = childContext.Value(subtotalAPIContextKey{}).(SubtotalAPI).
		GetByName(name)
	if e != nil {
		return
	}

	e = verify(assert.Equal,
		name,
		subtotal.Name(),
	)
	if e != nil {
		return
	}

	e = verify(assert.Equal,
		false,
		subtotal.HasParent(),
	)
	if e != nil {
		return
	}

	return
}

type (
	subtotalAPIContextKey struct{}
)

// TODO: make package

func verify(assertion assertionFunc, expected, actual interface{}) (e error) {
	var (
		t tee
	)

	assertion(&t, expected, actual)

	return t.e
}

type assertionFunc func(
	t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{},
) bool

type tee struct {
	e error
}

func (t *tee) Errorf(format string, args ...interface{}) {
	t.e = fmt.Errorf(format, args...)
}
