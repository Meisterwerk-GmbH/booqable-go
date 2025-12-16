package bq

import (
	"context"

	orderlist "github.com/Meisterwerk-GmbH/booqable-go/orders/list"
	"github.com/Meisterwerk-GmbH/booqable-go/orders/model"
)

// Type aliases so callers can keep importing bq.
type Order = model.Order
type OrderAttributes = model.OrderAttributes
type Planning = model.Planning
type PlanningAttributes = model.PlanningAttributes
type Customer = model.Customer
type CustomerAttributes = model.CustomerAttributes

// ListOrders returns all orders, including plannings and customers.
func (c *Client) ListOrders(ctx context.Context, opts ListOptions) ([]Order, error) {
	params := orderlist.Params{
		Page:  opts.Page,
		Per:   opts.Per,
		Query: opts.Query,
	}
	return orderlist.List(ctx, c, params)
}
