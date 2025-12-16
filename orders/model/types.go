package model

// Order represents an order with optional included resources.
type Order struct {
	ID         string          `json:"id"`
	Attributes OrderAttributes `json:"attributes"`
	Customer   *Customer       `json:"customer,omitempty"`
	Plannings  []Planning      `json:"plannings,omitempty"`
}

// OrderAttributes holds the key fields exposed by the Orders API.
type OrderAttributes struct {
	Number     string  `json:"number,omitempty"`
	Status     string  `json:"status,omitempty"`
	State      string  `json:"state,omitempty"`
	StartsAt   *string `json:"starts_at,omitempty"`
	EndsAt     *string `json:"ends_at,omitempty"`
	CustomerID *string `json:"customer_id,omitempty"`
	Reference  *string `json:"reference,omitempty"`
	CreatedAt  *string `json:"created_at,omitempty"`
	UpdatedAt  *string `json:"updated_at,omitempty"`
}

// Planning is the planned position on an order.
type Planning struct {
	ID         string             `json:"id"`
	Attributes PlanningAttributes `json:"attributes"`
}

// PlanningAttributes describes a planning item.
type PlanningAttributes struct {
	ProductID        *string `json:"product_id,omitempty"`
	ProductVariantID *string `json:"product_variant_id,omitempty"`
	ProductName      *string `json:"product_name,omitempty"`
	Quantity         *int    `json:"quantity,omitempty"`
	StartsAt         *string `json:"starts_at,omitempty"`
	EndsAt           *string `json:"ends_at,omitempty"`
}

// Customer represents the customer related to an order.
type Customer struct {
	ID         string             `json:"id"`
	Attributes CustomerAttributes `json:"attributes"`
}

// CustomerAttributes contains common customer details.
type CustomerAttributes struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	CompanyName *string `json:"company_name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
}
