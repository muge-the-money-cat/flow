package main

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

func TestSubtotal(t *testing.T) {
	var (
		testSuite = godog.TestSuite{
			ScenarioInitializer: initialiseSubtotalScenarios,
			Options:             testutils.GodogOptions,
		}
	)

	if testSuite.Run() != 0 {
		t.Fatal()
	}

	return
}

func initialiseSubtotalScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^there is a Subtotal API$`,
		thereIsASubtotalAPI,
	)
	ctx.Step(`^we POST a Subtotal with name "(.+)" and no parent$`,
		postASubtotalWithNoParent,
	)
	ctx.Step(`^there should be a Subtotal with name "(.+)" and no parent$`,
		thereShouldBeASubtotalWithNoParent,
	)
	ctx.Step(`^there is a Subtotal with name "(.+)"$`,
		thereIsASubtotal,
	)
	ctx.Step(`^we POST a Subtotal with name "(.+)" and parent "(.+)"$`,
		postASubtotalWithParent,
	)
	ctx.Step(`^there should be a Subtotal with name "(.+)" and parent "(.+)"$`,
		thereShouldBeASubtotalWithParent,
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
		nil,
		subtotal.Parent(),
	)
	if e != nil {
		return
	}

	return
}

func thereIsASubtotal(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	childContext, e = postASubtotalWithNoParent(parentContext, name)
	if e != nil {
		return
	}

	return
}

func postASubtotalWithParent(
	parentContext context.Context, name string, parentName string,
) (
	childContext context.Context, e error,
) {
	var (
		parent   Subtotal
		subtotal Subtotal
	)

	childContext = parentContext

	parent, e = childContext.Value(subtotalAPIContextKey{}).(SubtotalAPI).
		GetByName(parentName)
	if e != nil {
		return
	}

	subtotal = NewSubtotalWithParent(name, parent)

	e = childContext.Value(subtotalAPIContextKey{}).(SubtotalAPI).Post(subtotal)
	if e != nil {
		return
	}

	return
}

func thereShouldBeASubtotalWithParent(
	parentContext context.Context, name string, parentName string,
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
		parentName,
		subtotal.Parent().Name(),
	)
	if e != nil {
		return
	}

	return
}

type (
	subtotalAPIContextKey struct{}
)
