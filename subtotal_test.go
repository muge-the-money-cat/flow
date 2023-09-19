package main

import (
	"context"
	"net/http"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

var (
	subtotalURL string = testutils.EndpointURL(basePath, subtotalSubpath)
)

func initialiseSubtotalScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^a Subtotal endpoint is available$`,
		subtotalEndpointIsAvailable,
	)
	ctx.Step(`^we GET a Subtotal by name "(.+)"$`,
		getSubtotalByName,
	)
	ctx.Step(`^we POST a Subtotal with name "(.+)" and no parent$`,
		postSubtotalWithNoParent,
	)
	ctx.Step(`^we should see a Subtotal with name "(.+)" and no parent$`,
		shouldSeeSubtotalWithNoParent,
	)
	ctx.Step(`^we POST a Subtotal with name "(.+)" and parent "(.+)"$`,
		postSubtotalWithParent,
	)
	ctx.Step(`^we should see a Subtotal with name "(.+)" and parent "(.+)"$`,
		shouldSeeSubtotalWithParent,
	)
	ctx.Step(`^we PATCH a Subtotal named "(.+)" with new name "(.+)"$`,
		patchSubtotalWithNewName,
	)
	ctx.Step(`^we PATCH a Subtotal named "(.+)" with new parent "(.+)"$`,
		patchSubtotalWithNewParent,
	)
	ctx.Step(`^we DELETE a Subtotal named "(.+)"$`,
		deleteSubtotal,
	)

	return
}

func subtotalEndpointIsAvailable(parentContext context.Context) (
	childContext context.Context, e error,
) {
	var (
		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().Options(subtotalURL)
	if e != nil {
		return
	}

	e = testutils.Verify(assert.Equal,
		http.StatusNoContent,
		response.StatusCode(),
	)
	if e != nil {
		return
	}

	return
}

func getSubtotalByName(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		response *resty.Response
		subtotal Subtotal
	)

	childContext = parentContext

	response, subtotal, e = _getSubtotalByName(name)
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

func patchSubtotalWithNewName(parentContext context.Context,
	name, newName string,
) (
	childContext context.Context, e error,
) {
	var (
		response *resty.Response
		subtotal Subtotal
	)

	childContext = parentContext

	response, subtotal, e = _getSubtotalByName(name)
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

func patchSubtotalWithNewParent(parentContext context.Context,
	name, newParentName string,
) (
	childContext context.Context, e error,
) {
	var (
		response *resty.Response
		subtotal Subtotal
	)

	childContext = parentContext

	response, subtotal, e = _getSubtotalByName(name)
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

func _getSubtotalByName(name string) (
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
