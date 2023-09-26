package main

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow"
	"github.com/muge-the-money-cat/flow/testutils"
)

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
	ctx.Step(`^we should see error "(.+)"$`,
		shouldSeeError,
	)
	ctx.Step(`^we delete Subtotal "(.+)"$`,
		_deleteSubtotal,
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

	childContext, e = parseCLIOutput(parentContext)
	if e != nil {
		return
	}

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

func shouldSee(parentContext context.Context, expected string,
	actualContextKey any,
) (
	childContext context.Context, e error,
) {
	var (
		actual string = parentContext.Value(actualContextKey).(string)
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

func shouldSeeMessage(parentContext context.Context, expected string) (
	childContext context.Context, e error,
) {
	return shouldSee(parentContext, expected, cliOutputMessageContextKey{})
}

func shouldSeeError(parentContext context.Context, expected string) (
	childContext context.Context, e error,
) {
	return shouldSee(parentContext, expected, cliOutputErrorContextKey{})
}

func shouldSeeSubtotalWithParent(parentContext context.Context,
	name, parentName string,
) (
	childContext context.Context, e error,
) {
	var (
		expected = flow.Subtotal{
			Name:       name,
			ParentName: parentName,
		}

		actual flow.Subtotal = parentContext.
			Value(cliOutputPayloadContextKey{}).(flow.Subtotal)
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

func _deleteSubtotal(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		args = []string{appName,
			subtotalCommandName,
			subtotalDeleteCommandName,
			prefixFlag(subtotalNameFlag), name,
		}
	)

	childContext = parentContext

	e = run(args)
	if e != nil {
		return
	}

	childContext, e = parseCLIOutput(parentContext)
	if e != nil {
		return
	}

	return
}
