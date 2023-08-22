package testutils

import (
	"github.com/go-resty/resty/v2"
)

var (
	RESTClient *resty.Client = resty.New()
)
