package bq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ListOptions configures pagination and extra query params.
type ListOptions struct {
	Page  int
	Per   int
	Query url.Values
}

type listResponse struct {
	Data  json.RawMessage `json:"data"`
	Meta  *listMeta       `json:"meta,omitempty"`
	Links *listLinks      `json:"links,omitempty"`
}

type listMeta struct {
	Page        int `json:"page,omitempty"`
	CurrentPage int `json:"current_page,omitempty"`
	PerPage     int `json:"per_page,omitempty"`
	TotalPages  int `json:"total_pages,omitempty"`
}

type listLinks struct {
	Next string `json:"next,omitempty"`
}

// ListAll fetches every page for the given resource.
func (c *Client) ListAll(ctx context.Context, resource string, opts ListOptions) ([]json.RawMessage, error) {
	per := opts.Per
	if per <= 0 {
		per = 100
	}
	page := opts.Page
	if page <= 0 {
		page = 1
	}
	var all []json.RawMessage

	for {
		items, resp, err := c.list(ctx, resource, page, per, opts.Query)
		if err != nil {
			return nil, fmt.Errorf("list %s page %d: %w", resource, page, err)
		}
		all = append(all, items...)

		if !hasMore(resp, page, per, len(items)) {
			break
		}
		page++
	}
	return all, nil
}

func (c *Client) list(ctx context.Context, resource string, page, per int, q url.Values) ([]json.RawMessage, *listResponse, error) {
	params := cloneValues(q)
	params.Set("page", strconv.Itoa(page))
	params.Set("per", strconv.Itoa(per))

	var resp listResponse
	path := "/" + strings.TrimPrefix(resource, "/")
	if err := c.Get(ctx, path, params, &resp); err != nil {
		return nil, nil, err
	}

	items, err := decodeItems(resp.Data)
	if err != nil {
		return nil, nil, err
	}
	return items, &resp, nil
}

func hasMore(resp *listResponse, page, per, got int) bool {
	if resp == nil {
		return false
	}
	if resp.Meta != nil {
		current := resp.Meta.Page
		if current == 0 {
			current = resp.Meta.CurrentPage
		}
		if resp.Meta.TotalPages > 0 && current < resp.Meta.TotalPages {
			return true
		}
	}
	if resp.Links != nil && resp.Links.Next != "" {
		return true
	}
	return got == per && got > 0
}

func decodeItems(raw json.RawMessage) ([]json.RawMessage, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("decode items: %w", err)
	}
	return items, nil
}

func cloneValues(v url.Values) url.Values {
	if v == nil {
		return url.Values{}
	}
	out := make(url.Values, len(v))
	for key, vals := range v {
		cp := make([]string, len(vals))
		copy(cp, vals)
		out[key] = cp
	}
	return out
}
