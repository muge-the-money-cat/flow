package main

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

func TestPostSubtotalWithNoParent(t *testing.T) {
	var (
		testSuite = godog.TestSuite{
			ScenarioInitializer: initScenarioPostSubtotalWithNoParent,
			Options:             testutils.GodogOptions,
		}
	)

	if testSuite.Run() != 0 {
		t.Fatal()
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

	e = testutils.Verify(assert.Equal,
		name,
		subtotal.Name(),
	)
	if e != nil {
		return
	}

	e = testutils.Verify(assert.Equal,
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
