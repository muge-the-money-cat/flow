package main

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
	"github.com/muge-the-money-cat/flow/testutils"
)

const (
	basePath = "v1"
	root     = "/"
)

type flowHTTPAPIV1Server struct {
	entClient       *ent.Client
	ginEngine       *gin.Engine
	baseRouterGroup *gin.RouterGroup
}

func NewFlowHTTPAPIV1Server(entDriverName, entSourceName string,
	options ...flowHTTPAPIV1ServerOption,
) (
	a *flowHTTPAPIV1Server, e error,
) {
	a = new(flowHTTPAPIV1Server)

	e = a.initialiseEntClient(entDriverName, entSourceName)
	if e != nil {
		return
	}

	a.initialiseGinEngine()

	a.applyOptions(options)

	e = a.listenAndServe()
	if e != nil {
		return
	}

	return
}

func (a *flowHTTPAPIV1Server) initialiseEntClient(
	entDriverName, entSourceName string,
) (
	e error,
) {
	a.entClient, e = ent.Open(entDriverName, entSourceName)
	if e != nil {
		return
	}

	e = a.entClient.Schema.Create(
		context.Background(),
	)
	if e != nil {
		return
	}

	return
}

func (a *flowHTTPAPIV1Server) initialiseGinEngine() {
	a.ginEngine = gin.Default()

	a.baseRouterGroup = a.ginEngine.Group(basePath)

	return
}

func (a *flowHTTPAPIV1Server) applyOptions(
	options []flowHTTPAPIV1ServerOption,
) {
	var (
		option flowHTTPAPIV1ServerOption
	)

	for _, option = range options {
		option(a)
	}

	return
}

func (a *flowHTTPAPIV1Server) listenAndServe() (e error) {
	const (
		network = "tcp"
	)

	var (
		listener net.Listener
	)

	listener, e = net.Listen(network, testutils.TestServerAddress)
	if e != nil {
		return
	}

	go a.ginEngine.RunListener(listener)

	return
}

func (*flowHTTPAPIV1Server) handleError(c *gin.Context, e error) {
	switch {
	case ent.IsNotFound(e):
		c.Status(http.StatusNotFound)

	case ent.IsConstraintError(e):
		c.Status(http.StatusConflict)

	default:
		c.String(http.StatusInternalServerError,
			e.Error(), // XXX: remove before flight
		)
	}

	return
}

type flowHTTPAPIV1ServerOption func(*flowHTTPAPIV1Server)
