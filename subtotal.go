package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
	"github.com/muge-the-money-cat/flow/ent/subtotal"
)

const (
	subtotalNilID         = 0
	subtotalNilName       = ""
	subtotalNilParentName = ""
	subtotalSubpath       = "subtotal"
)

type Subtotal struct {
	ID         int
	Name       string
	ParentName string
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
		routerGroup.PATCH(root, server.patchSubtotal)

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
		q *ent.Subtotal
		s Subtotal

		create *ent.SubtotalCreate
	)

	ginContext.Bind(&s)

	create = server.entClient.Subtotal.Create().
		SetName(s.Name)

	if s.ParentName != subtotalNilParentName {
		q, e = server.getSubtotalByName(ginContext, s.ParentName)
		if e != nil {
			server.handleError(ginContext, e)

			return
		}

		create.SetParentID(q.ID)
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

	q, e = server.getSubtotalByName(ginContext, s.Name)
	if e != nil {
		server.handleError(ginContext, e)

		return
	}

	s = Subtotal{
		ID:   q.ID,
		Name: q.Name,
	}

	if q.Edges.Parent != nil {
		s.ParentName = q.Edges.Parent.Name
	}

	ginContext.JSON(http.StatusOK, s)

	return
}

func (server *flowHTTPAPIV1Server) patchSubtotal(ginContext *gin.Context) {
	var (
		e error
		q *ent.Subtotal
		s Subtotal

		update *ent.SubtotalUpdateOne
	)

	ginContext.Bind(&s)

	update = server.entClient.Subtotal.UpdateOneID(s.ID)

	if s.Name != subtotalNilName {
		update = update.SetName(s.Name)
	}

	if s.ParentName != subtotalNilParentName {
		q, e = server.getSubtotalByName(ginContext, s.ParentName)
		if e != nil {
			server.handleError(ginContext, e)

			return
		}

		update = update.SetParentID(q.ID)
	}

	_, e = update.Save(
		ginContext.Request.Context(),
	)
	if e != nil {
		server.handleError(ginContext, e)

		return
	}

	ginContext.Status(http.StatusNoContent)

	return
}

func (server *flowHTTPAPIV1Server) getSubtotalByName(ginContext *gin.Context,
	name string,
) (
	q *ent.Subtotal, e error,
) {
	q, e = server.entClient.Subtotal.Query().
		WithParent().
		Where(
			subtotal.Name(name),
		).
		Only(
			ginContext.Request.Context(),
		)
	if e != nil {
		return
	}

	return
}
