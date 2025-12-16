# bq-go

`bq-go` is a small Go client for the Booqable V4 REST API. Use it as a building block for CLIs or services.

## Install

```bash
go get github.com/Meisterwerk-GmbH/booqable-go
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
