package main

import (
	"io"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	"github.com/muge-the-money-cat/flow"
)

const (
	appName = "flow"
)

var (
	client *resty.Client  = resty.New()
	logger zerolog.Logger = zerolog.New(writer)
	writer io.Writer      = zerolog.ConsoleWriter{Out: os.Stderr}

	// FIXME: make configurable
	driverName    string = "sqlite3"
	serverAddress string = "127.78.88.89:8081"
	sourceName    string = "file:/tmp/flow?cache=shared&_fk=1"
	urlScheme     string = "http"
)

func main() {
	var (
		e error
	)

	defer die(&e)

	gin.SetMode(gin.ReleaseMode)

	_, e = flow.NewFlowV1HTTPAPIServer(
		serverAddress,
		driverName,
		sourceName,
		flow.WithSubtotalEndpoint(),
	)
	if e != nil {
		return
	}

	e = run(os.Args)
	if e != nil {
		return
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
		return
	}

	return
}

func die(pointer *error) {
	var (
		e error = *pointer
	)

	if e != nil {
		logger.Fatal().Err(e).Send()
	}

	return
}

func endpointURL(subpath string) string {
	var (
		base = &url.URL{
			Scheme: urlScheme,
			Host:   serverAddress,
			Path:   flow.BasePathV1,
		}
	)

	return base.JoinPath(subpath).String()
}

func prefixFlag(flag string) string {
	const (
		flagPrefix = "--"
	)

	return flagPrefix + flag
}
