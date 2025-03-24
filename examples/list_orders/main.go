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

	orders, err := client.Orders.List()
	if err != nil {
		log.Fatalf("List error: %v", err)
	}

	for i, order := range orders {
		fmt.Printf("Order %d: %v\n", i, order)
	}
}
