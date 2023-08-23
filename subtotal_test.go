package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

func TestSubtotal(t *testing.T) {
	var (
		testSuite = godog.TestSuite{
			ScenarioInitializer: initialiseSubtotalScenarios,
			Options:             testutils.GodogOptions,
		}
	)

	NewSubtotalHTTPAPIServer()

	if testSuite.Run() != 0 {
		t.Fatal()
	}

	return
}

func initialiseSubtotalScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^a Subtotal HTTP API server is up$`,
		subtotalHTTPAPIServerIsUp,
	)
	ctx.Step(`^we GET a Subtotal by name "(.+)"$`,
		getSubtotalByName,
	)
	ctx.Step(`^we should see HTTP response status (\d{3})$`,
		shouldSeeHTTPResponseStatus,
	)

	return
}

func subtotalHTTPAPIServerIsUp(parentContext context.Context) (
	childContext context.Context, e error,
) {
	const (
		urlFormat = "http://%s/up/"
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
		urlFormat = "http://%s/subtotal/?name=%s"
	)

	var (
		url string = fmt.Sprintf(urlFormat, testutils.TestServerAddress, name)

		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().Get(url)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		subtotalHTTPResponseContextKey{},
		response,
	)

	return
}

func shouldSeeHTTPResponseStatus(parentContext context.Context, status int) (
	childContext context.Context, e error,
) {
	childContext = parentContext

	e = testutils.Verify(assert.Equal,
		status,
		parentContext.Value(subtotalHTTPResponseContextKey{}).(*resty.Response).
			StatusCode(),
	)
	if e != nil {
		return
	}

	return
}

type (
	subtotalHTTPAPIContextKey      struct{}
	subtotalHTTPResponseContextKey struct{}
)
