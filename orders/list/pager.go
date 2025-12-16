package list

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Meisterwerk-GmbH/booqable-go/orders/model"
)

type pager struct {
	getter Getter
	query  url.Values
	per    int
	page   int
}

func newPager(getter Getter, params Params) pager {
	per := params.Per
	if per <= 0 {
		per = 100
	}
	page := params.Page
	if page <= 0 {
		page = 1
	}

	query := cloneValues(params.Query)
	ensureInclude(query, "customer")
	ensureInclude(query, "plannings")

	return pager{
		getter: getter,
		query:  query,
		per:    per,
		page:   page,
	}
}

func (p pager) fetchAll(ctx context.Context) ([]model.Order, error) {
	currentPage := p.page
	var orders []model.Order

	for {
		p.query.Set("page", strconv.Itoa(currentPage))
		p.query.Set("per", strconv.Itoa(p.per))

		var resp listResponse
		if err := p.getter.Get(ctx, "/orders", p.query, &resp); err != nil {
			return nil, fmt.Errorf("list orders page %d: %w", currentPage, err)
		}
		pageOrders, err := resp.decodeOrders()
		if err != nil {
			return nil, fmt.Errorf("decode orders page %d: %w", currentPage, err)
		}
		orders = append(orders, pageOrders...)

		if !hasMore(resp.Meta, resp.Links, currentPage, p.per, len(resp.Data)) {
			break
		}
		currentPage++
	}
	return orders, nil
}
