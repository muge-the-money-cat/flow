package main

import (
	"io"
	"log"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/muge-the-money-cat/flow"
	"github.com/urfave/cli/v2"
)

const (
	appName = "flow"
)

var (
	client *resty.Client = resty.New()
	buffer io.ReadWriter

	serverAddress string
)

func main() {
	var (
		e error
	)

	// TODO: set serverAddress

	buffer = os.Stdout

	e = run(os.Args)
	if e != nil {
		log.Fatalln(e)
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

	e = app.Run(args)
	if e != nil {
		return
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
