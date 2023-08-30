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

func flowHTTPAPIV1ServerIsUp(parentContext context.Context) (
	childContext context.Context, e error,
) {
	const (
		urlFormat = "http://%s/v1/up/"
	)

	var (
		url string = fmt.Sprintf(urlFormat, testutils.TestServerAddress)

		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().Get(url)
	if e != nil {
		return
	}

	e = testutils.Verify(assert.Equal,
		http.StatusOK,
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
	const (
		urlFormat = "http://%s/v1/subtotal/"
	)

	var (
		url string = fmt.Sprintf(urlFormat, testutils.TestServerAddress)

		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().
		SetQueryParam("Name", name).
		Get(url)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		subtotalHTTPResponseContextKey{},
		response,
	)

	return
}

func shouldSeeHTTPResponseStatus(parentContext context.Context, expected int) (
	childContext context.Context, e error,
) {
	var (
		actual int
	)

	childContext = parentContext

	actual = parentContext.
		Value(subtotalHTTPResponseContextKey{}).(*resty.Response).
		StatusCode()

	switch actual {
	case http.StatusBadRequest:
		fallthrough

	case http.StatusInternalServerError:
		e = fmt.Errorf(
			parentContext.
				Value(subtotalHTTPResponseContextKey{}).(*resty.Response).
				String(),
		)

		return
	}

	e = testutils.Verify(assert.Equal,
		expected,
		actual,
	)
	if e != nil {
		return
	}

	return
}

func postSubtotalWithNoParent(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	const (
		urlFormat = "http://%s/v1/subtotal/"
	)

	var (
		subtotal        = Subtotal{Name: name}
		url      string = fmt.Sprintf(urlFormat, testutils.TestServerAddress)

		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().
		SetBody(subtotal).
		SetResult(&subtotal).
		Post(url)
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
	childContext = parentContext

	e = testutils.Verify(assert.Equal,
		Subtotal{
			Name:     name,
			ParentID: 0,
		},
		parentContext.Value(subtotalHTTPResponseParsedContextKey{}).(Subtotal),
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
