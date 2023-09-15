package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

func endpointURL(host, subpath string) string {
	var (
		base = &url.URL{
			Scheme: "http",
			Host:   host,
			Path:   basePath,
		}
	)

	return base.JoinPath(subpath).String()
}

func shouldSeeHTTPResponseStatus(parentContext context.Context, expected int) (
	childContext context.Context, e error,
) {
	var (
		actual int = parentContext.Value(
			subtotalHTTPResponseContextKey{},
		).(*resty.Response).
			StatusCode()
	)

	childContext = parentContext

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
