package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/testutils"
)

type subtotalHTTPAPIServer struct {
	engine *gin.Engine
}

func NewSubtotalHTTPAPIServer() (s *subtotalHTTPAPIServer) {
	s = &subtotalHTTPAPIServer{
		engine: gin.Default(),
	}

	s.engine.GET("/up/", s.up)

	go s.engine.Run(testutils.TestServerAddress)

	return
}

func (s *subtotalHTTPAPIServer) up(c *gin.Context) {
	c.JSON(http.StatusOK, nil)

	return
}
