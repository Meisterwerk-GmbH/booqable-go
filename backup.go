package bq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

// BackupOptions configures a backup run.
type BackupOptions struct {
	Resources      []string
	OutputPath     string
	PerPage        int
	ResourceParams map[string]url.Values // map resource -> query params
}

// BackupAll fetches every resource and returns the raw JSON items.
func BackupAll(ctx context.Context, c *Client, opts BackupOptions) (map[string][]json.RawMessage, error) {
	if len(opts.Resources) == 0 {
		return nil, errors.New("no resources provided")
	}

	result := make(map[string][]json.RawMessage, len(opts.Resources))
	for _, resource := range opts.Resources {
		params := cloneValues(opts.ResourceParams[resource])
		items, err := c.ListAll(ctx, resource, ListOptions{
			Per:   opts.PerPage,
			Query: params,
		})
		if err != nil {
			return nil, fmt.Errorf("backup %s: %w", resource, err)
		}
		result[resource] = items
	}
	return result, nil
}

// BackupToFile writes the backup to disk.
func BackupToFile(ctx context.Context, c *Client, opts BackupOptions) (map[string][]json.RawMessage, error) {
	data, err := BackupAll(ctx, c, opts)
	if err != nil {
		return nil, err
	}
	if opts.OutputPath == "" {
		return data, nil
	}
	if err := writeJSON(opts.OutputPath, data); err != nil {
		return nil, err
	}
	return data, nil
}

func writeJSON(path string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create dirs: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
