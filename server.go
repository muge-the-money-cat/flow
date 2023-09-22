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

type flowV1HTTPAPIServer struct {
	entClient       *ent.Client
	ginEngine       *gin.Engine
	baseRouterGroup *gin.RouterGroup
}

func NewFlowV1HTTPAPIServer(address, entDriverName, entSourceName string,
	options ...flowV1HTTPAPIServerOption,
) (
	server *flowV1HTTPAPIServer, e error,
) {
	server = new(flowV1HTTPAPIServer)

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

func (server *flowV1HTTPAPIServer) initialiseEntClient(
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

func (server *flowV1HTTPAPIServer) initialiseGinEngine() {
	server.ginEngine = gin.Default()

	server.baseRouterGroup = server.ginEngine.Group(basePath)

	return
}

func (server *flowV1HTTPAPIServer) applyOptions(
	options []flowV1HTTPAPIServerOption,
) {
	var (
		option flowV1HTTPAPIServerOption
	)

	for _, option = range options {
		option(server)
	}

	return
}

func (server *flowV1HTTPAPIServer) listenAndServe(address string) (e error) {
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

func (*flowV1HTTPAPIServer) handleError(c *gin.Context, pointer *error) {
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

type flowV1HTTPAPIServerOption func(*flowV1HTTPAPIServer)
