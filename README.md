# Todos
- there have to be cleaned-up many things, it is kind of a mess right now
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
orders, err := client.ListOrders(ctx, bq.ListOptions{Per: 100})
if err != nil {
    log.Fatal(err)
}

for _, order := range orders {
    fmt.Println("Order", order.Attributes.Number)
    if order.Customer != nil {
        fmt.Println("Customer", order.Customer.Attributes.FirstName, order.Customer.Attributes.LastName)
    }
    for _, planning := range order.Plannings {
        fmt.Println("Planning", planning.Attributes.ProductName, planning.Attributes.Quantity)
    }
}
```

## Notes

- Authentication uses the `X-Api-Key` header.
- Pagination uses `page`/`per` query params. The client will also honor `links.next` or `meta.total_pages` if present.

### Generic listing

`ListAll` remains available if you need to retrieve other resources without the higher level helpers.
