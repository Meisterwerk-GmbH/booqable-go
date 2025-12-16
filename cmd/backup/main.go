package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	bq "github.com/booqable/bq-go"
)

func main() {
	var (
		baseURL   string
		token     string
		resources string
		outPath   string
		per       int
		timeout   time.Duration
	)

	flag.StringVar(&baseURL, "base-url", bq.DefaultBaseURL, "Booqable API base URL")
	flag.StringVar(&token, "token", "", "Booqable API token (defaults to BQ_API_TOKEN)")
	flag.StringVar(&resources, "resources", "products,customers,orders,order_items,stock_items,bookings", "comma-separated resource names to backup")
	flag.StringVar(&outPath, "out", "backup.json", "path to write the backup JSON")
	flag.IntVar(&per, "per", 100, "items per page")
	flag.DurationVar(&timeout, "timeout", 60*time.Second, "request timeout")
	flag.Parse()

	if token == "" {
		token = os.Getenv("BQ_API_TOKEN")
	}
	if token == "" {
		log.Fatal("missing API token: use -token or set BQ_API_TOKEN")
	}

	resourceList := splitAndClean(resources)
	if len(resourceList) == 0 {
		log.Fatal("no resources provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := bq.NewClient(token, bq.WithBaseURL(baseURL), bq.WithHTTPClient(&http.Client{Timeout: timeout}))
	backupOpts := bq.BackupOptions{
		Resources:  resourceList,
		OutputPath: outPath,
		PerPage:    per,
	}

	result, err := bq.BackupToFile(ctx, client, backupOpts)
	if err != nil {
		log.Fatalf("backup failed: %v", err)
	}
	log.Printf("backup complete: %d resources written to %s", len(result), outPath)
}

func splitAndClean(input string) []string {
	parts := strings.Split(input, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
