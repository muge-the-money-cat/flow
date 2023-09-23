package main

import (
	"bytes"
	"testing"

	"github.com/cucumber/godog"

	"github.com/muge-the-money-cat/flow"
	"github.com/muge-the-money-cat/flow/testutils"
)

const (
	testServerAddress = "127.78.88.89:8081"
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

	serverAddress = testServerAddress

	_, e = flow.NewFlowV1HTTPAPIServer(
		serverAddress,
		driverName,
		sourceName,
		flow.WithSubtotalEndpoint(),
	)
	if e != nil {
		t.Fatal(e)
	}

	buffer = new(bytes.Buffer)

	if testSuite.Run() != 0 {
		t.Fatal()
	}

	return
}

func initialiseScenarios(ctx *godog.ScenarioContext) {
	initialiseSubtotalScenarios(ctx)

	return
}
