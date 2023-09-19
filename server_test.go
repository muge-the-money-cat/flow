package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

type (
	httpResponseContextKey       struct{}
	httpResponseParsedContextKey struct{}
)

func TestFlowHTTPAPIV1Server(t *testing.T) {
	const (
		entDriverName = "sqlite3"
		entSourceName = "file:ent?mode=memory&cache=shared&_fk=1"
	)

	var (
		godogOptions = &godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
		}
		testSuite = godog.TestSuite{
			ScenarioInitializer: initialiseScenarios,
			Options:             godogOptions,
		}

		e error
	)

	_, e = NewFlowHTTPAPIV1Server(
		testutils.TestServerAddress,
		entDriverName,
		entSourceName,
		withSubtotalEndpoint(),
		withAccountEndpoint(),
	)
	if e != nil {
		log.Fatalln(e)
	}

	if testSuite.Run() != 0 {
		t.Fatal()
	}

	return
}

func initialiseScenarios(ctx *godog.ScenarioContext) {
	initialiseAccountScenarios(ctx)
	initialiseGenericScenarios(ctx)
	initialiseSubtotalScenarios(ctx)

	return
}

func initialiseGenericScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^we should see HTTP response status (\d{3})$`,
		shouldSeeHTTPResponseStatus,
	)

	return
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
