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

type subtotalCLIOutput struct {
	Message  string
	Subtotal flow.Subtotal
}

func initialiseSubtotalScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^we create Subtotal "(.+)" with no parent$`,
		createSubtotalWithNoParent,
	)
	ctx.Step(`^we should see message "(.+)"$`,
		shouldSeeMessage,
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
			subtotalCommandName,
			subtotalCreateCommandName,
			prefixFlag(subtotalNameFlag), name,
		}

		output      subtotalCLIOutput
		outputBytes []byte
	)

	childContext = parentContext

	e = run(args)
	if e != nil {
		return
	}

	outputBytes, e = io.ReadAll(buffer)
	if e != nil {
		return
	}

	e = json.Unmarshal(outputBytes, &output)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		cliOutputMessageContextKey{},
		output.Message,
	)

	childContext = context.WithValue(childContext,
		cliOutputPayloadContextKey{},
		output.Subtotal,
	)

	return
}

func shouldSeeMessage(parentContext context.Context, expected string) (
	childContext context.Context, e error,
) {
	var (
		actual string = parentContext.
			Value(cliOutputMessageContextKey{}).(string)
	)

	childContext = parentContext

	e = testutils.Verify(assert.Equal,
		expected,
		actual,
	)
	if e != nil {
		return
	}

	return
}

func shouldSeeSubtotalWithNoParent(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		actual flow.Subtotal = parentContext.
			Value(cliOutputPayloadContextKey{}).(flow.Subtotal)

		expected = flow.Subtotal{
			Name: name,
		}
	)

	childContext = parentContext

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
