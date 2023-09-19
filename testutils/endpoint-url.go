package testutils

import (
	"net/url"
)

func EndpointURL(basepath, subpath string) string {
	var (
		base = &url.URL{
			Scheme: "http",
			Host:   TestServerAddress,
			Path:   basepath,
		}
	)

	return base.JoinPath(subpath).String()
}
