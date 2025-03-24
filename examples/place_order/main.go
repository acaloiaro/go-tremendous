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

	campaigns, err := client.Campaigns.List()
	if err != nil {
		fmt.Println("unable to list campaigns", err)
		return
	}

	newOrderArgs := tremendous.OrderArgs{
		Reward: tremendous.RewardArg{
			CampaignID: &campaigns[0].ID,
			Recipient: tremendous.OrderRecipient{
				Name:  "Testy McTesterson",
				Email: "testy@example.com",
			},
			Value: tremendous.OrderDenomination{
				Denomination: 1.00,
				CurrencyCode: "USD",
			},
			Delivery: tremendous.OrderDelivery{
				Method: "LINK",
			},
		},
		Payment: tremendous.OrderPaymentArg{
			FundingSourceID: "BALANCE",
		},
	}
	// Example: Create an order (replace with valid fields).
	newOrder, err := client.Orders.Create(newOrderArgs)
	if err != nil {
		log.Fatalf("Create error: %v", err)
	}
	fmt.Println("Created order:", newOrder, "\n\n")

	// Example: Retrieve that order by ID.
	retOrder, err := client.Orders.Retrieve(newOrder.ID)
	if err != nil {
		log.Fatalf("Retrieve error: %v", err)
	}
	fmt.Println("Retrieved order:", retOrder)
}
