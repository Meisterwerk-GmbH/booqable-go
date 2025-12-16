# bq-go

`bq-go` is a small Go client for the Booqable V4 REST API. The initial focus is a `backup` CLI that exports data from common resources to a single JSON file for local safekeeping.

## Quick start

```bash
cd bq-go
# set your API token (from Booqable dashboard)
export BQ_API_TOKEN=your_token_here
go run ./cmd/backup \
  -resources products,customers,orders,order_items,stock_items,bookings \
  -out ./backup.json
```

Flags:

- `-token` API token. Defaults to `BQ_API_TOKEN` environment variable.
- `-base-url` API base URL (default `https://api.booqable.com/v4`).
- `-resources` comma-separated resources to fetch.
- `-out` path to the JSON backup file.
- `-per` items per page (default 100).
- `-timeout` request timeout (default 60s).

The output JSON groups pages by resource:

```json
{
  "products": [ { /* product */ }, ... ],
  "customers": [ { /* customer */ }, ... ]
}
```

## Library usage

```go
client := bq.NewClient("BQ_API_TOKEN")
ctx := context.Background()
items, err := client.ListAll(ctx, "products", bq.ListOptions{PerPage: 100})
```

## Notes

- Authentication uses the `X-Api-Key` header.
- Pagination uses `page`/`per` query params. The client will also honor `links.next` or `meta.total_pages` if present.
- The resource names used by the CLI are configurable; add or remove as needed for your account.
