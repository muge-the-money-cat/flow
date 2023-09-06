package main

import (
	"context"
	"fmt"
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
	ctx.Step(`^a Flow HTTP API v1 server is up$`,
		flowHTTPAPIV1ServerIsUp,
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

		actual Subtotal
	)

	childContext = parentContext

	actual = parentContext.Value(
		subtotalHTTPResponseParsedContextKey{},
	).(Subtotal)

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
