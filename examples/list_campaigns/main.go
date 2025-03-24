package main

import (
	"cmp"
	"fmt"
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
	campaigns, err := client.Campaigns.List()
	if err != nil {
		fmt.Printf("unable to list campaigns: %s\n", err)
		return
	}

	for _, c := range campaigns {
		fmt.Println(c)
	}
}
