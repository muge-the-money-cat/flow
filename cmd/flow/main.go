package main

import (
	"io"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/muge-the-money-cat/flow"
)

const (
	appName = "flow"
)

var (
	client *resty.Client = resty.New()

	logger zerolog.Logger
	writer io.Writer

	serverAddress string
)

func main() {
	var (
		e error
	)

	// TODO: set serverAddress

	writer = os.Stderr

	e = run(os.Args)
	if e != nil {
		log.Fatal().Err(e).Send()
	}

	return
}

func run(args []string) (e error) {
	var (
		app = &cli.App{
			Name:  appName,
			Usage: "accounting \u00d7 zen",
			Commands: []*cli.Command{
				subtotalCommand,
			},
		}
	)

	logger = zerolog.New(writer)

	e = app.Run(args)
	if e != nil {
		logger.Fatal().Err(e).Send()
	}

	return
}

func endpointURL(subpath string) string {
	var (
		base = &url.URL{
			Scheme: "http", // FIXME
			Host:   serverAddress,
			Path:   flow.BasePathV1,
		}
	)

	return base.JoinPath(subpath).String()
}
