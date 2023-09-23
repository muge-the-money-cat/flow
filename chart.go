package flow

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muge-the-money-cat/flow/ent"
)

const (
	chartSubpath                = "chart"
	chartQueryParamSubtotalName = "Root"
)

type Chart struct {
	Edges []ChartEdge
}

func (c *Chart) appendEdge(tail, head string) {
	c.Edges = append(c.Edges,
		ChartEdge{
			Tail: tail,
			Head: head,
		},
	)

	return
}

type ChartEdge struct {
	Tail string
	Head string
	// > In the edge (x, y) directed from x to y,
	// > the vertices x and y are called the endpoints of the edge,
	// > x the tail of the edge and y the head of the edge.
	// (https://en.wikipedia.org/wiki/Graph_theory#Directed_graph)
}

func withChartEndpoint() (option flowV1HTTPAPIServerOption) {
	option = func(server *flowV1HTTPAPIServer) {
		var (
			routerGroup *gin.RouterGroup = server.baseRouterGroup.Group(
				chartSubpath,
			)
		)

		routerGroup.OPTIONS(root, server.chartOptions)

		routerGroup.GET(root, server.getChart)

		return
	}

	return
}

func (server *flowV1HTTPAPIServer) chartOptions(ginContext *gin.Context) {
	ginContext.Status(http.StatusNoContent) // FIXME

	return
}

func (server *flowV1HTTPAPIServer) getChart(ginContext *gin.Context) {
	var (
		c Chart
		e error
	)

	defer server.handleError(ginContext, &e)

	e = server.populateChart(&c,
		ginContext.Query(chartQueryParamSubtotalName),
		ginContext.Request.Context(),
	)
	if e != nil {
		return
	}

	ginContext.JSON(http.StatusOK, c)

	return
}

func (server *flowV1HTTPAPIServer) populateChart(chart *Chart,
	subtotalName string, ctx context.Context,
) (
	e error,
) {
	var (
		a *ent.Account
		q *ent.Subtotal
		r *ent.Subtotal
	)

	q, e = server.getSubtotalByName(subtotalName,
		ctx,
		loadSubtotalChildren(),
		loadSubtotalAccounts(),
	)
	if e != nil {
		return
	}

	for _, r = range q.Edges.Children {
		chart.appendEdge(q.Name, r.Name)

		server.populateChart(chart, r.Name, ctx)
		// recursive depth-first traversal
	}

	for _, a = range q.Edges.Accounts {
		chart.appendEdge(q.Name, a.Name)
	}

	return
}
