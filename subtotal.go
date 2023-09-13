package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
	"github.com/muge-the-money-cat/flow/ent/subtotal"
)

const (
	subtotalSubpath = "subtotal"
)

type Subtotal struct {
	Name     string
	ParentID int
}

func withSubtotalEndpoint() (option flowHTTPAPIV1ServerOption) {
	option = func(server *flowHTTPAPIV1Server) {
		var (
			routerGroup *gin.RouterGroup = server.baseRouterGroup.Group(
				subtotalSubpath,
			)
		)

		routerGroup.OPTIONS(root, server.subtotalOptions)
		routerGroup.POST(root, server.postSubtotal)
		routerGroup.GET(root, server.getSubtotal)

		return
	}

	return
}

func (server *flowHTTPAPIV1Server) subtotalOptions(ginContext *gin.Context) {
	ginContext.Status(http.StatusNoContent) // FIXME

	return
}

func (server *flowHTTPAPIV1Server) postSubtotal(ginContext *gin.Context) {
	var (
		e error
		s Subtotal

		create *ent.SubtotalCreate
	)

	ginContext.Bind(&s)

	create = server.entClient.Subtotal.Create().
		SetName(s.Name)

	if s.ParentID != 0 {
		create.SetParentID(s.ParentID)
	}

	_, e = create.Save(
		ginContext.Request.Context(),
	)
	if e != nil {
		server.handleError(ginContext, e)

		return
	}

	ginContext.Status(http.StatusCreated)

	return
}

func (server *flowHTTPAPIV1Server) getSubtotal(ginContext *gin.Context) {
	var (
		e error
		q *ent.Subtotal
		s Subtotal
	)

	ginContext.Bind(&s)

	q, e = server.entClient.Subtotal.Query().
		WithParent().
		Where(
			subtotal.Name(s.Name),
		).
		Only(
			ginContext.Request.Context(),
		)
	if e != nil {
		server.handleError(ginContext, e)

		return
	}

	s = Subtotal{
		Name: q.Name,
	}

	if q.Edges.Parent != nil {
		s.ParentID = q.Edges.Parent.ID
	}

	ginContext.JSON(http.StatusOK, s)

	return
}
