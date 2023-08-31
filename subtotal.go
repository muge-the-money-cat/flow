package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
	"github.com/muge-the-money-cat/flow/ent/subtotal"
)

type Subtotal struct {
	Name     string `json:"name"`
	ParentID int    `json:"parentID"`
}

func (a *flowHTTPAPIV1Server) postSubtotal(c *gin.Context) {
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
		if ent.IsConstraintError(e) {
			c.Status(http.StatusConflict)

			return
		}

		c.String(http.StatusInternalServerError,
			e.Error(), // XXX: remove before flight
		)

		return
	}

	c.Status(http.StatusCreated)

	return
}

func (a *flowHTTPAPIV1Server) getSubtotal(c *gin.Context) {
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
