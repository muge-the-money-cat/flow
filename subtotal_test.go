package main

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

var (
	subtotalURL string = testutils.EndpointURL(basePath, subtotalSubpath)
)

func initialiseSubtotalScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^Subtotal endpoint is available$`,
		newEndpointAvailableScenarioStep(subtotalURL),
	)
	ctx.Step(`^we GET Subtotal "(.+)"$`,
		getSubtotal,
	)
	ctx.Step(`^we POST Subtotal "(.+)" with no parent$`,
		postSubtotalWithNoParent,
	)
	ctx.Step(`^we should see Subtotal "(.+)" with no parent$`,
		shouldSeeSubtotalWithNoParent,
	)
	ctx.Step(`^we POST Subtotal "(.+)" with parent "(.+)"$`,
		postSubtotalWithParent,
	)
	ctx.Step(`^we should see Subtotal "(.+)" with parent "(.+)"$`,
		shouldSeeSubtotalWithParent,
	)
	ctx.Step(`^we PATCH Subtotal "(.+)" with new name "(.+)"$`,
		patchSubtotalName,
	)
	ctx.Step(`^we PATCH Subtotal "(.+)" with new parent "(.+)"$`,
		patchSubtotalParent,
	)
	ctx.Step(`^we DELETE Subtotal "(.+)"$`,
		deleteSubtotal,
	)

	return
}

func getSubtotal(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		response *resty.Response
		subtotal Subtotal
	)

	childContext = parentContext

	response, subtotal, e = getSubtotalByName(name)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	childContext = context.WithValue(childContext,
		httpResponseParsedContextKey{},
		subtotal,
	)

	return
}

func postSubtotalWithParent(parentContext context.Context,
	name, parentName string,
) (
	childContext context.Context, e error,
) {
	var (
		subtotal = Subtotal{
			Name:       name,
			ParentName: parentName,
		}

		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().
		SetBody(subtotal).
		Post(subtotalURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	return
}

func postSubtotalWithNoParent(parentContext context.Context, name string) (
	context.Context, error,
) {
	return postSubtotalWithParent(parentContext, name, nilSubtotalParentName)
}

func shouldSeeSubtotalWithParent(parentContext context.Context,
	name, parentName string,
) (
	childContext context.Context, e error,
) {
	var (
		expected = Subtotal{
			Name:       name,
			ParentName: parentName,
		}

		actual Subtotal = parentContext.Value(
			httpResponseParsedContextKey{},
		).(Subtotal)
	)

	childContext = parentContext

	actual.ID = nilSubtotalID

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
	context.Context, error,
) {
	return shouldSeeSubtotalWithParent(parentContext,
		name,
		nilSubtotalParentName,
	)
}

func patchSubtotalName(parentContext context.Context,
	name, newName string,
) (
	childContext context.Context, e error,
) {
	var (
		response *resty.Response
		subtotal Subtotal
	)

	childContext = parentContext

	_, subtotal, e = getSubtotalByName(name)
	if e != nil {
		return
	}

	subtotal = Subtotal{
		ID:   subtotal.ID,
		Name: newName,
	}

	response, e = testutils.RESTClient.R().
		SetBody(subtotal).
		Patch(subtotalURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	return
}

func patchSubtotalParent(parentContext context.Context,
	name, newParentName string,
) (
	childContext context.Context, e error,
) {
	var (
		response *resty.Response
		subtotal Subtotal
	)

	childContext = parentContext

	_, subtotal, e = getSubtotalByName(name)
	if e != nil {
		return
	}

	subtotal = Subtotal{
		ID:         subtotal.ID,
		ParentName: newParentName,
	}

	response, e = testutils.RESTClient.R().
		SetBody(subtotal).
		Patch(subtotalURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	return
}

func deleteSubtotal(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		response *resty.Response
		subtotal Subtotal
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().
		SetQueryParam("Name", name).
		SetResult(&subtotal).
		Delete(subtotalURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	childContext = context.WithValue(childContext,
		httpResponseParsedContextKey{},
		subtotal,
	)

	return
}

func getSubtotalByName(name string) (
	response *resty.Response, subtotal Subtotal, e error,
) {
	response, e = testutils.RESTClient.R().
		SetQueryParam("Name", name).
		SetResult(&subtotal).
		Get(subtotalURL)
	if e != nil {
		return
	}

	return
}
