package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
	"github.com/muge-the-money-cat/flow/ent/subtotal"
	"github.com/muge-the-money-cat/flow/testutils"
)

type Subtotal struct {
	Name     string `json:"name"`
	ParentID int    `json:"parentID"`
}

type subtotalHTTPAPIServer struct {
	entClient *ent.Client
	ginEngine *gin.Engine
}

func NewSubtotalHTTPAPIServer(entDriverName, entSourceName string) (
	a *subtotalHTTPAPIServer, e error,
) {
	a = &subtotalHTTPAPIServer{
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

	a.ginEngine.GET("/up/", a.up)
	a.ginEngine.POST("/subtotal/", a.post)
	a.ginEngine.GET("/subtotal/", a.get)

	go a.ginEngine.Run(testutils.TestServerAddress)

	return
}

func (a *subtotalHTTPAPIServer) up(c *gin.Context) {
	c.Status(http.StatusOK)

	return
}

func (a *subtotalHTTPAPIServer) post(c *gin.Context) {
	var (
		e error
		s Subtotal

		create *ent.SubtotalCreate
	)

	c.Bind(&s)

	create = a.entClient.Subtotal.Create().
		SetName(s.Name)

	if s.ParentID != 0 {
		create.SetParentID(s.ParentID)
	}

	_, e = create.Save(
		c.Request.Context(),
	)
	if e != nil {
		c.Status(http.StatusInternalServerError)

		return
	}

	c.Status(http.StatusCreated)

	return
}

func (a *subtotalHTTPAPIServer) get(c *gin.Context) {
	var (
		e error
		q *ent.Subtotal
		s Subtotal
	)

	c.Bind(&s)

	q, e = a.entClient.Subtotal.Query().
		WithParent().
		Where(
			subtotal.Name(s.Name),
		).
		Only(
			c.Request.Context(),
		)
	if e != nil {
		if ent.IsNotFound(e) {
			c.Status(http.StatusNotFound)

			return
		}

		c.Status(http.StatusInternalServerError)

		return
	}

	s = Subtotal{
		Name: q.Name,
	}

	if q.Edges.Parent != nil {
		s.ParentID = q.Edges.Parent.ID
	}

	c.JSON(http.StatusOK, s)

	return
}
