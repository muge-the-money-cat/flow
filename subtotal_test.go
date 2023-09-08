package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
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

func postSubtotalWithNoParent(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		subtotal = Subtotal{Name: name}

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

func shouldSeeSubtotalWithNoParent(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		expected = Subtotal{
			Name:     name,
			ParentID: 0,
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

type (
	subtotalHTTPAPIContextKey            struct{}
	subtotalHTTPResponseContextKey       struct{}
	subtotalHTTPResponseParsedContextKey struct{}
)

const (
	subtotalURLFormat = "http://%s/v1/subtotal"
)

var (
	subtotalURL string = fmt.Sprintf(subtotalURLFormat,
		testutils.TestServerAddress,
	)
)
