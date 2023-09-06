package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

func flowHTTPAPIV1ServerIsUp(parentContext context.Context) (
	childContext context.Context, e error,
) {
	const (
		urlFormat = "http://%s/v1/up"
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

func shouldSeeHTTPResponseStatus(parentContext context.Context, expected int) (
	childContext context.Context, e error,
) {
	var (
		actual int
	)

	childContext = parentContext

	actual = parentContext.Value(
		subtotalHTTPResponseContextKey{},
	).(*resty.Response).
		StatusCode()

	switch actual {
	case http.StatusBadRequest:
		fallthrough

	case http.StatusInternalServerError:
		e = fmt.Errorf(
			parentContext.Value(
				subtotalHTTPResponseContextKey{},
			).(*resty.Response).
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
