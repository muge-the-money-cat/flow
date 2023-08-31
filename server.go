package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
	"github.com/muge-the-money-cat/flow/testutils"
)

const (
	root = "/"
)

type flowHTTPAPIV1Server struct {
	entClient *ent.Client
	ginEngine *gin.Engine
}

func NewFlowHTTPAPIV1Server(entDriverName, entSourceName string) (
	a *flowHTTPAPIV1Server, e error,
) {
	var (
		subtotal *gin.RouterGroup
		up       *gin.RouterGroup
		v1       *gin.RouterGroup
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

	v1 = a.ginEngine.Group("/v1")

	up = v1.Group("/up")

	up.GET(root, a.up)

	subtotal = v1.Group("/subtotal")

	subtotal.POST(root, a.postSubtotal)
	subtotal.GET(root, a.getSubtotal)

	go a.ginEngine.Run(testutils.TestServerAddress)

	return
}

func (a *flowHTTPAPIV1Server) up(c *gin.Context) {
	c.Status(http.StatusOK)

	return
}
