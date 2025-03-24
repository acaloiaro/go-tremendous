package tremendous

import (
	"fmt"
	"net/http"
)

// TestflightURL is the testflight API endpoint.
const TestflightURL = "https://testflight.tremendous.com/api/v2"

// ProductionURL is the production API endpoint
const ProductionURL = "https://api.tremendous.com/api/v2"

// Client is the primary entry point for using the Tremendous API
type Client struct {
	APIKey    string
	BaseURL   string
	Campaigns *CampaignsService
	Orders    *OrdersService
	Products  *ProductsService
}

type ClientArgs struct {
	ApiKey     string // the api key used to authenticate to the API
	Production bool   // whether to use the production or testflight API
}

// NewClient constructs a new Tremendous client using the given API key.
func NewClient(args ClientArgs) *Client {
	baseURL := TestflightURL
	if args.Production {
		baseURL = ProductionURL
	}
	client := &Client{
		APIKey:  args.ApiKey,
		BaseURL: baseURL,
	}
	client.Campaigns = &CampaignsService{client: client}
	client.Orders = &OrdersService{client: client}
	client.Products = &ProductsService{client: client}
	return client
}

// headers returns the default headers for requests, including authorization.
func (c *Client) headers() http.Header {
	h := http.Header{}
	h.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "application/json")
	h.Set("User-Agent", "go-tremendous/1.0.0")
	return h
}
