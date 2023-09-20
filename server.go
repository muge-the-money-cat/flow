package main

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
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

func NewFlowHTTPAPIV1Server(address, entDriverName, entSourceName string,
	options ...flowHTTPAPIV1ServerOption,
) (
	server *flowHTTPAPIV1Server, e error,
) {
	server = new(flowHTTPAPIV1Server)

	e = server.initialiseEntClient(entDriverName, entSourceName)
	if e != nil {
		return
	}

	server.initialiseGinEngine()

	server.applyOptions(options)

	e = server.listenAndServe(address)
	if e != nil {
		return
	}

	return
}

func (server *flowHTTPAPIV1Server) initialiseEntClient(
	entDriverName, entSourceName string,
) (
	e error,
) {
	server.entClient, e = ent.Open(entDriverName, entSourceName)
	if e != nil {
		return
	}

	e = server.entClient.Schema.Create(
		context.Background(),
	)
	if e != nil {
		return
	}

	return
}

func (server *flowHTTPAPIV1Server) initialiseGinEngine() {
	server.ginEngine = gin.Default()

	server.baseRouterGroup = server.ginEngine.Group(basePath)

	return
}

func (server *flowHTTPAPIV1Server) applyOptions(
	options []flowHTTPAPIV1ServerOption,
) {
	var (
		option flowHTTPAPIV1ServerOption
	)

	for _, option = range options {
		option(server)
	}

	return
}

func (server *flowHTTPAPIV1Server) listenAndServe(address string) (e error) {
	const (
		network = "tcp"
	)

	var (
		listener net.Listener
	)

	listener, e = net.Listen(network, address)
	if e != nil {
		return
	}

	go server.ginEngine.RunListener(listener)

	return
}

func (*flowHTTPAPIV1Server) handleError(c *gin.Context, pointer *error) {
	var (
		e error = *pointer
		// NOTE: This DEFERRED function takes a pointer argument because
		// > A deferred function's arguments are evaluated
		// > when the defer statement is evaluated.
		// (https://go.dev/blog/defer-panic-and-recover)
	)

	switch {
	case e == nil:
		return

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
