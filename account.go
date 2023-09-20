package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
	"github.com/muge-the-money-cat/flow/ent/account"
)

const (
	accountSubpath         = "account"
	nilAccountID           = 0
	nilAccountName         = ""
	nilAccountSubtotalName = ""
)

type Account struct {
	ID           int
	Name         string
	SubtotalName string
}

func newAccountFromEntAccount(q *ent.Account) (a Account) {
	a = Account{
		ID:           q.ID,
		Name:         q.Name,
		SubtotalName: q.Edges.Subtotal.Name,
	}

	return
}

func withAccountEndpoint() (option flowHTTPAPIV1ServerOption) {
	option = func(server *flowHTTPAPIV1Server) {
		var (
			routerGroup *gin.RouterGroup = server.baseRouterGroup.Group(
				accountSubpath,
			)
		)

		routerGroup.OPTIONS(root, server.accountOptions)

		routerGroup.DELETE(root, server.deleteAccount)
		routerGroup.GET(root, server.getAccount)
		routerGroup.PATCH(root, server.patchAccount)
		routerGroup.POST(root, server.postAccount)

		return
	}

	return
}

func (server *flowHTTPAPIV1Server) accountOptions(ginContext *gin.Context) {
	ginContext.Status(http.StatusNoContent) // FIXME

	return
}

func (server *flowHTTPAPIV1Server) postAccount(ginContext *gin.Context) {
	var (
		a Account
		e error
		s *ent.Subtotal
	)

	defer server.handleError(ginContext, &e)

	ginContext.Bind(&a)

	s, e = server.getSubtotalByName(a.SubtotalName,
		ginContext.Request.Context(),
	)
	if e != nil {
		return
	}

	_, e = server.entClient.Account.Create().
		SetName(a.Name).
		SetSubtotalID(s.ID).
		Save(
			ginContext.Request.Context(),
		)
	if e != nil {
		return
	}

	ginContext.Status(http.StatusNoContent)

	return
}

func (server *flowHTTPAPIV1Server) getAccount(ginContext *gin.Context) {
	var (
		a Account
		e error
		q *ent.Account
	)

	defer server.handleError(ginContext, &e)

	ginContext.Bind(&a)

	q, e = server.getAccountByName(a.Name,
		ginContext.Request.Context(),
	)
	if e != nil {
		return
	}

	a = newAccountFromEntAccount(q)

	ginContext.JSON(http.StatusOK, a)

	return
}

func (server *flowHTTPAPIV1Server) patchAccount(ginContext *gin.Context) {
	var (
		a Account
		e error
		s *ent.Subtotal

		update *ent.AccountUpdateOne
	)

	defer server.handleError(ginContext, &e)

	ginContext.Bind(&a)

	update = server.entClient.Account.UpdateOneID(a.ID)

	if a.Name != nilAccountName {
		update = update.SetName(a.Name)
	}

	if a.SubtotalName != nilAccountSubtotalName {
		s, e = server.getSubtotalByName(a.SubtotalName,
			ginContext.Request.Context(),
		)
		if e != nil {
			return
		}

		update = update.SetSubtotalID(s.ID)
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

func (server *flowHTTPAPIV1Server) deleteAccount(ginContext *gin.Context) {
	var (
		a Account
		e error
		q *ent.Account
	)

	defer server.handleError(ginContext, &e)

	ginContext.Bind(&a)

	q, e = server.getAccountByName(a.Name,
		ginContext.Request.Context(),
	)
	if e != nil {
		return
	}

	e = server.entClient.Account.DeleteOneID(q.ID).
		Exec(
			ginContext.Request.Context(),
		)
	if e != nil {
		return
	}

	a = newAccountFromEntAccount(q)

	ginContext.JSON(http.StatusOK, a)

	return
}

func (server *flowHTTPAPIV1Server) getAccountByName(name string,
	ctx context.Context,
) (
	q *ent.Account, e error,
) {
	q, e = server.entClient.Account.Query().
		WithSubtotal().
		Where(
			account.Name(name),
		).
		Only(ctx)
	if e != nil {
		return
	}

	return
}
