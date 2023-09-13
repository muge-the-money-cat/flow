package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

var (
	subtotalURL string = baseURL.
		JoinPath(subtotalSubpath).
		String()
)

func TestSubtotal(t *testing.T) {
	var (
		testSuite = godog.TestSuite{
			ScenarioInitializer: initialiseSubtotalScenarios,
			Options:             testutils.GodogOptions,
		}

		e error
	)

	_, e = NewFlowHTTPAPIV1Server(
		testutils.EntDriverName,
		testutils.EntSourceName,
		withSubtotalEndpoint(),
	)
	if e != nil {
		t.Fatal(e)
	}

	if testSuite.Run() != 0 {
		t.Fatal()
	}

	return
}

func initialiseSubtotalScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^a Subtotal endpoint is available$`,
		subtotalEndpointIsAvailable,
	)
	ctx.Step(`^we GET a Subtotal by name "(.+)"$`,
		getSubtotalByName,
	)
	ctx.Step(`^we should see HTTP response status (\d{3})$`,
		shouldSeeHTTPResponseStatus,
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
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().
		SetQueryParam("Name", name).
		Get(subtotalURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		subtotalHTTPResponseContextKey{},
		response,
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
		SetResult(&subtotal).
		Post(subtotalURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		subtotalHTTPResponseContextKey{},
		response,
	)

	childContext = context.WithValue(childContext,
		subtotalHTTPResponseParsedContextKey{},
		subtotal,
	)

	return
}

func postSubtotalWithNoParent(parentContext context.Context, name string) (
	context.Context, error,
) {
	return postSubtotalWithParent(parentContext, name, nilParentName)
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
			subtotalHTTPResponseParsedContextKey{},
		).(Subtotal)
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
	context.Context, error,
) {
	return shouldSeeSubtotalWithParent(parentContext, name, nilParentName)
}

type (
	subtotalHTTPAPIContextKey            struct{}
	subtotalHTTPResponseContextKey       struct{}
	subtotalHTTPResponseParsedContextKey struct{}
)
