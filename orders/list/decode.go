package list

import (
	"encoding/json"
	"fmt"

	"github.com/Meisterwerk-GmbH/booqable-go/orders/model"
)

type listResponse struct {
	Data     []orderResource    `json:"data"`
	Included []includedResource `json:"included,omitempty"`
	Meta     *listMeta          `json:"meta,omitempty"`
	Links    *listLinks         `json:"links,omitempty"`
}

type orderResource struct {
	ID            string                `json:"id"`
	Type          string                `json:"type"`
	Attributes    model.OrderAttributes `json:"attributes"`
	Relationships orderRelationships    `json:"relationships"`
}

type orderRelationships struct {
	Plannings relMany `json:"plannings,omitempty"`
	Customer  relOne  `json:"customer,omitempty"`
}

type relMany struct {
	Data []relationRef `json:"data"`
}

type relOne struct {
	Data *relationRef `json:"data"`
}

type relationRef struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type includedResource struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes json.RawMessage `json:"attributes"`
}

func (r listResponse) decodeOrders() ([]model.Order, error) {
	includedByType := make(map[string]map[string]includedResource)
	for _, inc := range r.Included {
		if _, ok := includedByType[inc.Type]; !ok {
			includedByType[inc.Type] = make(map[string]includedResource)
		}
		includedByType[inc.Type][inc.ID] = inc
	}

	var orders []model.Order
	for _, item := range r.Data {
		order := model.Order{
			ID:         item.ID,
			Attributes: item.Attributes,
		}

		if rel := item.Relationships.Customer.Data; rel != nil {
			order.Customer = decodeCustomer(rel.ID, includedByType["customers"])
		}
		for _, rel := range item.Relationships.Plannings.Data {
			planning, err := decodePlanning(rel.ID, includedByType["plannings"])
			if err != nil {
				return nil, fmt.Errorf("decode planning %s: %w", rel.ID, err)
			}
			order.Plannings = append(order.Plannings, planning)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func decodePlanning(id string, included map[string]includedResource) (model.Planning, error) {
	if inc, ok := included[id]; ok {
		var attrs model.PlanningAttributes
		if len(inc.Attributes) > 0 {
			if err := json.Unmarshal(inc.Attributes, &attrs); err != nil {
				return model.Planning{}, err
			}
		}
		return model.Planning{
			ID:         inc.ID,
			Attributes: attrs,
		}, nil
	}
	return model.Planning{ID: id}, nil
}

func decodeCustomer(id string, included map[string]includedResource) *model.Customer {
	if inc, ok := included[id]; ok {
		var attrs model.CustomerAttributes
		if len(inc.Attributes) > 0 {
			if err := json.Unmarshal(inc.Attributes, &attrs); err != nil {
				return &model.Customer{ID: inc.ID}
			}
		}
		return &model.Customer{
			ID:         inc.ID,
			Attributes: attrs,
		}
	}
	return &model.Customer{ID: id}
}
