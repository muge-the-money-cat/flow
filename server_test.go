package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

const (
	testServerAddress = "127.78.88.89:8080"
)

func TestMain(m *testing.M) {
	const (
		entDriverName = "sqlite3"
		entSourceName = "file:ent?mode=memory&cache=shared&_fk=1"
	)

	var (
		e        error
		exitCode int
	)

	_, e = NewFlowHTTPAPIV1Server(
		testServerAddress,
		entDriverName,
		entSourceName,
		withSubtotalEndpoint(),
		withAccountEndpoint(),
	)
	if e != nil {
		log.Fatalln(e)
	}

	exitCode = m.Run()

	os.Exit(exitCode)
}

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
			httpResponseContextKey{},
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
				httpResponseContextKey{},
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

type (
	httpResponseContextKey       struct{}
	httpResponseParsedContextKey struct{}
)
