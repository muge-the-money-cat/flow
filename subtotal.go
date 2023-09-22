package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
	"github.com/muge-the-money-cat/flow/ent/subtotal"
)

const (
	nilSubtotalID         = 0
	nilSubtotalName       = ""
	nilSubtotalParentName = ""
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

func withSubtotalEndpoint() (option flowV1HTTPAPIServerOption) {
	option = func(server *flowV1HTTPAPIServer) {
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

func (server *flowV1HTTPAPIServer) subtotalOptions(ginContext *gin.Context) {
	ginContext.Status(http.StatusNoContent) // FIXME

	return
}

func (server *flowV1HTTPAPIServer) postSubtotal(ginContext *gin.Context) {
	var (
		e error
		q *ent.Subtotal
		s Subtotal

		create *ent.SubtotalCreate
	)

	defer server.handleError(ginContext, &e)

	ginContext.Bind(&s)

	create = server.entClient.Subtotal.Create().
		SetName(s.Name)

	if s.ParentName != nilSubtotalParentName {
		q, e = server.getSubtotalByName(s.ParentName,
			ginContext.Request.Context(),
		)
		if e != nil {
			return
		}

		create.SetParentID(q.ID)
	}

	_, e = create.Save(
		ginContext.Request.Context(),
	)
	if e != nil {
		return
	}

	ginContext.Status(http.StatusNoContent)

	return
}

func (server *flowV1HTTPAPIServer) getSubtotal(ginContext *gin.Context) {
	var (
		e error
		q *ent.Subtotal
		s Subtotal
	)

	defer server.handleError(ginContext, &e)

	ginContext.Bind(&s)

	q, e = server.getSubtotalByName(s.Name,
		ginContext.Request.Context(),
	)
	if e != nil {
		return
	}

	s = newSubtotalFromEntSubtotal(q)

	ginContext.JSON(http.StatusOK, s)

	return
}

func (server *flowV1HTTPAPIServer) patchSubtotal(ginContext *gin.Context) {
	var (
		e error
		q *ent.Subtotal
		s Subtotal

		update *ent.SubtotalUpdateOne
	)

	defer server.handleError(ginContext, &e)

	ginContext.Bind(&s)

	update = server.entClient.Subtotal.UpdateOneID(s.ID)

	if s.Name != nilSubtotalName {
		update = update.SetName(s.Name)
	}

	if s.ParentName != nilSubtotalParentName {
		q, e = server.getSubtotalByName(s.ParentName,
			ginContext.Request.Context(),
		)
		if e != nil {
			return
		}

		update = update.SetParentID(q.ID)
	}

	_, e = update.Save(
		ginContext.Request.Context(),
	)
	if e != nil {
		return
	}

	ginContext.Status(http.StatusNoContent)

	return
}

func (server *flowV1HTTPAPIServer) deleteSubtotal(ginContext *gin.Context) {
	var (
		e error
		q *ent.Subtotal
		s Subtotal
	)

	defer server.handleError(ginContext, &e)

	ginContext.Bind(&s)

	q, e = server.getSubtotalByName(s.Name,
		ginContext.Request.Context(),
		loadSubtotalChildren(),
	)
	if e != nil {
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
		return
	}

	s = newSubtotalFromEntSubtotal(q)

	ginContext.JSON(http.StatusOK, s)

	return
}

func (server *flowV1HTTPAPIServer) getSubtotalByName(name string,
	ctx context.Context, options ...subtotalQueryOption,
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

	q, e = query.Only(ctx)
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

func loadSubtotalAccounts() (option subtotalQueryOption) {
	option = func(query *ent.SubtotalQuery) {
		query = query.WithAccounts()

		return
	}

	return
}
