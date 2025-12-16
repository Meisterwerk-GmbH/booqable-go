package list

import (
	"context"
	"fmt"

	"github.com/Meisterwerk-GmbH/booqable-go/orders/model"
)

// List returns all orders, including their plannings and customers.
func List(ctx context.Context, getter Getter, params Params) ([]model.Order, error) {
	pager := newPager(getter, params)
	orders, err := pager.fetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	return orders, nil
}
