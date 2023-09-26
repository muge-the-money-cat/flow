package main

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli/v2"

	"github.com/muge-the-money-cat/flow"
)

const (
	subtotalCommandName       = "subtotal"
	subtotalCreateCommandName = "create"
	subtotalLogKey            = "subtotal"
	subtotalNameFlag          = "name"
	subtotalParentNameFlag    = "parent"
)

var (
	subtotalCommand *cli.Command = command(subtotalCommandName,
		subcommand(subtotalCreateCommandName, createSubtotal,
			stringFlag(subtotalNameFlag, true),
			stringFlag(subtotalParentNameFlag, false),
		),
	)
)

func createSubtotal(c *cli.Context) (e error) {
	var (
		subtotalErrorConflict error = fmt.Errorf(
			"Subtotal with same name exists",
		)
		subtotalErrorNotFound error = fmt.Errorf(
			"Parent Subtotal does not exist",
		)

		subtotal = flow.Subtotal{
			Name:       c.String(subtotalNameFlag),
			ParentName: c.String(subtotalParentNameFlag),
		}

		response *resty.Response
	)

	response, e = client.R().
		SetBody(subtotal).
		Post(
			subtotalURL(),
		)
	if e != nil {
		return
	}

	switch response.StatusCode() {
	case http.StatusNotFound:
		logger.Error().
			Err(subtotalErrorNotFound).
			Send()

		return

	case http.StatusConflict:
		_, subtotal, e = getSubtotalByName(subtotal.Name)
		if e != nil {
			return
		}

		logger.Error().
			Err(subtotalErrorConflict).
			Interface(subtotalLogKey, subtotal).
			Send()

		return

		// TODO: more cases
	}

	response, subtotal, e = getSubtotalByName(subtotal.Name)
	if e != nil {
		return
	}

	if response.StatusCode() != http.StatusOK {
		// TODO

		return
	}

	logger.Info().
		Interface(subtotalLogKey, subtotal).
		Msg("Subtotal successfully created")

	return
}

func getSubtotalByName(name string) (
	response *resty.Response, subtotal flow.Subtotal, e error,
) {
	response, e = client.R().
		SetQueryParam(flow.SubtotalQueryParamName, name).
		SetResult(&subtotal).
		Get(
			subtotalURL(),
		)
	if e != nil {
		return
	}

	return
}

func subtotalURL() string {
	return endpointURL(flow.SubtotalSubpath)
}
