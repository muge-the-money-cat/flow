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
	root = "/"
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
	const (
		basePath = "v1"
		network  = "tcp"
	)

	var (
		listener net.Listener
		option   flowHTTPAPIV1ServerOption
	)

	a = &flowHTTPAPIV1Server{
		ginEngine: gin.Default(),
	}

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

	a.baseRouterGroup = a.ginEngine.Group(basePath)

	for _, option = range options {
		option(a)
	}

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
