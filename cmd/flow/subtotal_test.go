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
	ctx.Step(`^we create Subtotal "(.+)" with parent "(.+)"$`,
		createSubtotalWithParent,
	)
	ctx.Step(`^we should see Subtotal "(.+)" with parent "(.+)"$`,
		shouldSeeSubtotalWithParent,
	)

	return
}

func createSubtotalWithParent(parentContext context.Context,
	name, parentName string,
) (
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

	if parentName != flow.NilSubtotalParentName {
		args = append(args,
			prefixFlag(subtotalParentNameFlag),
			parentName,
		)
	}

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

func createSubtotalWithNoParent(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	childContext, e = createSubtotalWithParent(parentContext,
		name,
		flow.NilSubtotalParentName,
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

func shouldSeeSubtotalWithParent(parentContext context.Context,
	name, parentName string,
) (
	childContext context.Context, e error,
) {
	var (
		actual flow.Subtotal = parentContext.
			Value(cliOutputPayloadContextKey{}).(flow.Subtotal)

		expected = flow.Subtotal{
			Name:       name,
			ParentName: parentName,
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

func shouldSeeSubtotalWithNoParent(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	childContext, e = shouldSeeSubtotalWithParent(parentContext,
		name,
		flow.NilSubtotalParentName,
	)

	return
}
