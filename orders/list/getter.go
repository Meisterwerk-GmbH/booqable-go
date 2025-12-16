package list

import (
	"context"
	"net/url"
)

// Getter represents the subset of Client used for listing.
type Getter interface {
	Get(ctx context.Context, p string, query url.Values, v interface{}) error
}
