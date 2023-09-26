package main

import (
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

	if response.StatusCode() != http.StatusNoContent {
		// TODO

		return
	}

	response, e = client.R().
		SetQueryParam(flow.SubtotalQueryParamName, subtotal.Name).
		SetResult(&subtotal).
		Get(
			subtotalURL(),
		)
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

func subtotalURL() string {
	return endpointURL(flow.SubtotalSubpath)
}
