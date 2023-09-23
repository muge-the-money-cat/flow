package main

import (
	"context"
	"encoding/json"
	"io"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow"
	"github.com/muge-the-money-cat/flow/testutils"
)

type (
	cliOutputContextKey struct{}
)

func initialiseSubtotalScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^we create Subtotal "(.+)" with no parent$`,
		createSubtotalWithNoParent,
	)
	ctx.Step(`^we should see Subtotal "(.+)" with no parent$`,
		shouldSeeSubtotalWithNoParent,
	)

	return
}

func createSubtotalWithNoParent(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		args = []string{appName,
			"subtotal",
			"create",
			"--name", name,
		}

		output []byte
	)

	childContext = parentContext

	e = run(args)
	if e != nil {
		return
	}

	output, e = io.ReadAll(buffer)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		cliOutputContextKey{},
		output,
	)

	return
}

func shouldSeeSubtotalWithNoParent(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		actual   flow.Subtotal
		expected = flow.Subtotal{
			Name: name,
		}

		output []byte = parentContext.Value(cliOutputContextKey{}).([]byte)
	)

	childContext = parentContext

	e = json.Unmarshal(output, &actual)

	actual.ID = flow.NilSubtotalID

	e = testutils.Verify(assert.Equal,
		expected,
		actual,
	)
	if e != nil {
		return
	}

	return
}
