package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/muge-the-money-cat/flow"
	"github.com/muge-the-money-cat/flow/testutils"
)

var (
	buffer *bytes.Buffer
)

type (
	cliOutputErrorContextKey   struct{}
	cliOutputMessageContextKey struct{}
	cliOutputPayloadContextKey struct{}
)

func TestFlowV1CLI(t *testing.T) {
	const (
		driverName = "sqlite3"
		sourceName = "file:cli?mode=memory&cache=shared&_fk=1"
	)

	var (
		testSuite = godog.TestSuite{
			ScenarioInitializer: initialiseScenarios,
			Options:             testutils.GodogOptions,
		}

		e error
	)

	_, e = flow.NewFlowV1HTTPAPIServer(gin.Default,
		serverAddress,
		driverName,
		sourceName,
		flow.WithSubtotalEndpoint(),
	)
	if e != nil {
		t.Fatal(e)
	}

	buffer = new(bytes.Buffer)

	writer = io.MultiWriter(
		zerolog.ConsoleWriter{Out: os.Stderr},
		buffer,
	)

	if testSuite.Run() != 0 {
		t.Fatal()
	}

	return
}

func initialiseScenarios(ctx *godog.ScenarioContext) {
	initialiseSubtotalScenarios(ctx)

	return
}

type cliOutput struct {
	Error    string
	Message  string
	Subtotal flow.Subtotal
}

func parseCLIOutput(parentContext context.Context) (
	childContext context.Context, e error,
) {
	var (
		b []byte
		o cliOutput
	)

	b, e = io.ReadAll(buffer)
	if e != nil {
		return
	}

	e = json.Unmarshal(b, &o)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		cliOutputErrorContextKey{},
		o.Error,
	)

	childContext = context.WithValue(childContext,
		cliOutputMessageContextKey{},
		o.Message,
	)

	childContext = context.WithValue(childContext,
		cliOutputPayloadContextKey{},
		o.Subtotal,
	)

	return
}
