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

func newSubtotalFromEntSubtotal(q *ent.Subtotal) (s Subtotal) {
	s = Subtotal{
		ID:   q.ID,
		Name: q.Name,
	}

	if q.Edges.Parent != nil {
		s.ParentName = q.Edges.Parent.Name
	}

	return
}

func withSubtotalEndpoint() (option flowHTTPAPIV1ServerOption) {
	option = func(server *flowHTTPAPIV1Server) {
		var (
			routerGroup *gin.RouterGroup = server.baseRouterGroup.Group(
				subtotalSubpath,
			)
		)

		routerGroup.OPTIONS(root, server.subtotalOptions)

		routerGroup.DELETE(root, server.deleteSubtotal)
		routerGroup.GET(root, server.getSubtotal)
		routerGroup.PATCH(root, server.patchSubtotal)
		routerGroup.POST(root, server.postSubtotal)

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

	ginContext.Status(http.StatusNoContent)

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

	s = newSubtotalFromEntSubtotal(q)

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

func (server *flowHTTPAPIV1Server) deleteSubtotal(ginContext *gin.Context) {
	var (
		e error
		q *ent.Subtotal
		s Subtotal
	)

	ginContext.Bind(&s)

	q, e = server.getSubtotalByName(ginContext, s.Name,
		loadSubtotalChildren(),
	)
	if e != nil {
		server.handleError(ginContext, e)

		return
	}

	if len(q.Edges.Children) > 0 {
		ginContext.Status(http.StatusConflict)

		return
	}

	e = server.entClient.Subtotal.DeleteOneID(q.ID).
		Exec(
			ginContext.Request.Context(),
		)
	if e != nil {
		server.handleError(ginContext, e)

		return
	}

	s = newSubtotalFromEntSubtotal(q)

	ginContext.JSON(http.StatusOK, s)

	return
}

func (server *flowHTTPAPIV1Server) getSubtotalByName(ginContext *gin.Context,
	name string, options ...subtotalQueryOption,
) (
	q *ent.Subtotal, e error,
) {
	var (
		option subtotalQueryOption
		query  *ent.SubtotalQuery
	)

	query = server.entClient.Subtotal.Query().
		WithParent().
		Where(
			subtotal.Name(name),
		)

	for _, option = range options {
		option(query)
	}

	q, e = query.Only(
		ginContext.Request.Context(),
	)
	if e != nil {
		return
	}

	return
}

type subtotalQueryOption func(*ent.SubtotalQuery)

func loadSubtotalChildren() (option subtotalQueryOption) {
	option = func(query *ent.SubtotalQuery) {
		query = query.WithChildren()

		return
	}

	return
}
