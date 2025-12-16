# Todos
- the client is to generic, it should have functions that return actual typed resources.
The most important resource at the Moment is the order-resources and it should be possible to retrieve all resources
including their positions and customers.
- extract multiple things into submodules to create a better overview

# booqable-go

`booqable-go` is a small Go client for the Booqable V4 REST API. Use it as a building block for CLIs or services.

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
