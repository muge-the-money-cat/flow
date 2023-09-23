package flow

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

const (
	accountQueryParamName = "Name"
)

var (
	accountURL string = testutils.EndpointURL(BasePathV1, accountSubpath)
)

func initialiseAccountScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^Account endpoint is available$`,
		newEndpointAvailableScenarioStep(accountURL),
	)
	ctx.Step(`^we GET Account "(.+)"$`,
		getAccount,
	)
	ctx.Step(`^we POST Account "(.+)" with Subtotal "(.+)"$`,
		postAccount,
	)
	ctx.Step(`^we should see Account "(.+)" with Subtotal "(.+)"$`,
		shouldSeeAccount,
	)
	ctx.Step(`^we PATCH Account "(.+)" with new name "(.+)"$`,
		patchAccountName,
	)
	ctx.Step(`^we PATCH Account "(.+)" with new Subtotal "(.+)"$`,
		patchAccountSubtotal,
	)
	ctx.Step(`^we DELETE Account "(.+)"$`,
		deleteAccount,
	)

	return
}

func getAccount(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		account  Account
		response *resty.Response
	)

	childContext = parentContext

	response, account, e = getAccountByName(name)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	childContext = context.WithValue(childContext,
		httpResponseParsedContextKey{},
		account,
	)

	return
}

func postAccount(parentContext context.Context, name, subtotalName string) (
	childContext context.Context, e error,
) {
	var (
		account = Account{
			Name:         name,
			SubtotalName: subtotalName,
		}

		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().
		SetBody(account).
		Post(accountURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	return
}

func shouldSeeAccount(parentContext context.Context,
	name, subtotalName string,
) (
	childContext context.Context, e error,
) {
	var (
		expected = Account{
			Name:         name,
			SubtotalName: subtotalName,
		}

		actual Account = parentContext.Value(
			httpResponseParsedContextKey{},
		).(Account)
	)

	childContext = parentContext

	actual.ID = nilAccountID

	e = testutils.Verify(assert.Equal,
		expected,
		actual,
	)
	if e != nil {
		return
	}

	return
}

func patchAccountName(parentContext context.Context,
	name, newName string,
) (
	childContext context.Context, e error,
) {
	var (
		account  Account
		response *resty.Response
	)

	childContext = parentContext

	_, account, e = getAccountByName(name)
	if e != nil {
		return
	}

	account = Account{
		ID:   account.ID,
		Name: newName,
	}

	response, e = testutils.RESTClient.R().
		SetBody(account).
		Patch(accountURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	return
}

func patchAccountSubtotal(parentContext context.Context,
	name, newSubtotalName string,
) (
	childContext context.Context, e error,
) {
	var (
		account  Account
		response *resty.Response
	)

	childContext = parentContext

	_, account, e = getAccountByName(name)
	if e != nil {
		return
	}

	account = Account{
		ID:           account.ID,
		SubtotalName: newSubtotalName,
	}

	response, e = testutils.RESTClient.R().
		SetBody(account).
		Patch(accountURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	return
}

func deleteAccount(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		account  Account
		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().
		SetQueryParam(accountQueryParamName, name).
		SetResult(&account).
		Delete(accountURL)
	if e != nil {
		return
	}

	childContext = context.WithValue(parentContext,
		httpResponseContextKey{},
		response,
	)

	childContext = context.WithValue(childContext,
		httpResponseParsedContextKey{},
		account,
	)

	return
}

func getAccountByName(name string) (
	response *resty.Response, account Account, e error,
) {
	response, e = testutils.RESTClient.R().
		SetQueryParam(accountQueryParamName, name).
		SetResult(&account).
		Get(accountURL)
	if e != nil {
		return
	}

	return
}
