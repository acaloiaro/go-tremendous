package main

import (
	"cmp"
	"fmt"
	"log"
	"os"

	"github.com/acaloiaro/go-tremendous"
)

func main() {
	// Replace with your actual API key.
	args := tremendous.ClientArgs{
		ApiKey:     cmp.Or(os.Getenv("TREMENDOUS_API_KEY"), "YOUR_API_KEY"),
		Production: false,
	}
	client := tremendous.NewClient(args)

	products, err := client.Products.List(&tremendous.ListProductsOptions{
		Country: "US", // list only US products
	})
	if err != nil {
		log.Fatalf("List error: %v", err)
	}

	for i, product := range products {
		fmt.Printf("Product %d: %v\n", i, product)
	}
}
