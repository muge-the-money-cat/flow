package main

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

var (
	chartURL string = testutils.EndpointURL(basePath, chartSubpath)
)

func initialiseChartScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^Chart endpoint is available$`,
		newEndpointAvailableScenarioStep(chartURL),
	)
	ctx.Step(`^we GET Chart based on Subtotal "(.+)"$`,
		getChart,
	)
	ctx.Step(`^we should see Chart with edge "(.+)" -> "(.+)"$`,
		shouldSeeChartWithEdge,
	)

	return
}

func getChart(parentContext context.Context, subtotalName string) (
	childContext context.Context, e error,
) {
	var (
		chart    Chart
		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().
		SetQueryParam(chartQueryParamSubtotalName, subtotalName).
		SetResult(&chart).
		Get(chartURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	childContext = context.WithValue(childContext,
		httpResponseParsedContextKey{},
		chart,
	)

	return
}

func shouldSeeChartWithEdge(parentContext context.Context, from, to string) (
	childContext context.Context, e error,
) {
	var (
		actual Chart = parentContext.Value(
			httpResponseParsedContextKey{},
		).(Chart)
	)

	childContext = parentContext

	e = testutils.Verify(assert.Contains,
		actual.Edges,
		ChartEdge{
			Tail: from,
			Head: to,
		},
	)
	if e != nil {
		return
	}

	return
}
