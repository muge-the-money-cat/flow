package main

import (
	"context"
	"net/http"

	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/muge-the-money-cat/flow/testutils"
)

var (
	accountURL string = testutils.EndpointURL(basePath, accountSubpath)
)

func initialiseAccountScenarios(ctx *godog.ScenarioContext) {
	ctx.Step(`^Account endpoint is available$`,
		accountEndpointIsAvailable,
	)
	ctx.Step(`^we GET Account "(.+)"$`,
		getAccountByName,
	)
	ctx.Step(`^we POST Account "(.+)" with Subtotal "(.+)"$`,
		postAccount,
	)
	ctx.Step(`^we should see Account "(.+)" with Subtotal "(.+)"$`,
		shouldSeeAccount,
	)
	ctx.Step(`^we PATCH Account "(.+)" with new name "(.+)"$`,
		patchAccountWithNewName,
	)
	ctx.Step(`^we PATCH Account "(.+)" with new Subtotal "(.+)"$`,
		patchAccountWithNewSubtotal,
	)
	ctx.Step(`^we DELETE Account "(.+)"$`,
		deleteAccount,
	)

	return
}

func accountEndpointIsAvailable(parentContext context.Context) (
	childContext context.Context, e error,
) {
	var (
		response *resty.Response
	)

	childContext = parentContext

	response, e = testutils.RESTClient.R().Options(accountURL)
	if e != nil {
		return
	}

	e = testutils.Verify(assert.Equal,
		http.StatusNoContent,
		response.StatusCode(),
	)
	if e != nil {
		return
	}

	return
}

func getAccountByName(parentContext context.Context, name string) (
	childContext context.Context, e error,
) {
	var (
		account  Account
		response *resty.Response
	)

	childContext = parentContext

	response, account, e = _getAccountByName(name)
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

func patchAccountWithNewName(parentContext context.Context,
	name, newName string,
) (
	childContext context.Context, e error,
) {
	var (
		account  Account
		response *resty.Response
	)

	childContext = parentContext

	_, account, e = _getAccountByName(name)
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

func patchAccountWithNewSubtotal(parentContext context.Context,
	name, newSubtotalName string,
) (
	childContext context.Context, e error,
) {
	var (
		account  Account
		response *resty.Response
	)

	childContext = parentContext

	_, account, e = _getAccountByName(name)
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
		SetQueryParam("Name", name).
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

func _getAccountByName(name string) (
	response *resty.Response, account Account, e error,
) {
	response, e = testutils.RESTClient.R().
		SetQueryParam("Name", name).
		SetResult(&account).
		Get(accountURL)
	if e != nil {
		return
	}

	return
}
