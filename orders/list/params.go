package list

import "net/url"

// Params captures pagination and query configuration for listing orders.
type Params struct {
	Page  int
	Per   int
	Query url.Values
}
