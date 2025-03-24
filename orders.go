package tremendous

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// OrdersService handles the "orders" resource.
type OrdersService struct {
	client *Client
}

// Order represents the JSON structure of an order in the Tremendous API.
type Order struct {
	ID         string    `json:"id,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	Message    string    `json:"message,omitempty"`
	Status     string    `json:"status,omitempty"`
	Payment    Payment   `json:"payment,omitempty"`
	Rewards    []Reward  `json:"rewards"`
}

// Payment captures details of the payment object used by the order.
type Payment struct {
	FundingSourceID string  `json:"funding_source_id,omitempty"`
	Amount          float64 `json:"amount,omitempty"`
	CurrencyCode    string  `json:"currency_code,omitempty"`
	Object          string  `json:"object,omitempty"`
}

// Reward captures details of a reward that has been paid out
type Reward struct {
	ID        string            `json:"id"`
	OrderID   string            `json:"order_id"`
	CreatedAt time.Time         `json:"created_at"`
	Value     OrderDenomination `json:"value"`
	Delivery  OrderDelivery     `json:"delivery"`
	Recipient OrderRecipient    `json:"recipient"`
}

// OrderArgs is the request struct for creating an order. It wraps
// the "order" object exactly as Tremendous expects in the JSON payload.
type OrderArgs struct {
	CampaignID *string         `json:"campaign_id"`
	ExternalID *string         `json:"external_id"`
	Payment    OrderPaymentArg `json:"payment"`
	Reward     RewardArg       `json:"reward"`
}

// OrderPaymentArg describes an order's payment details
type OrderPaymentArg struct {
	FundingSourceID string `json:"funding_source_id,omitempty"`
}

// OrderDenomination denominates the reward's monetary value
type OrderDenomination struct {
	Denomination float64 `json:"denomination"`
	CurrencyCode string  `json:"currency_code"`
}

// OrderDelivery designates how an order's reward is delivered
type OrderDelivery struct {
	Method string `json:"method"` // email, link, or phone
	Status string `json:"status"`
	Link   string `json:"link"`
}

// RewardArg is subset of field for on order reward
type RewardArg struct {
	CampaignID *string           `json:"campaign_id"`
	Delivery   OrderDelivery     `json:"delivery"`
	Products   []string          `json:"products"`
	Recipient  OrderRecipient    `json:"recipient"`
	Value      OrderDenomination `json:"value"`
}

// OrderRecipient contains details about reward recipients
type OrderRecipient struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// singleOrderResponse models the JSON for a single order return.
type singleOrderResponse struct {
	Order Order `json:"order"`
}

// listOrdersResponse models the JSON for a multi-order listing.
type listOrdersResponse struct {
	Orders []Order `json:"orders"`
}

// List fetches all orders. (Query params, pagination, etc., may be added later.)
func (s *OrdersService) List() ([]Order, error) {
	url := fmt.Sprintf("%s/orders", s.client.BaseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = s.client.headers()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var ordersResp listOrdersResponse
	if err := json.NewDecoder(resp.Body).Decode(&ordersResp); err != nil {
		return nil, fmt.Errorf("failed to decode orders list: %w", err)
	}
	return ordersResp.Orders, nil
}

// Retrieve fetches a single order by its ID.
func (s *OrdersService) Retrieve(orderID string) (*Order, error) {
	orderID = strings.TrimSpace(orderID)
	url := fmt.Sprintf("%s/orders/%s", s.client.BaseURL, orderID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = s.client.headers()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var singleResp singleOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&singleResp); err != nil {
		return nil, fmt.Errorf("failed to decode single order: %w", err)
	}
	return &singleResp.Order, nil
}

// Create sends a new order to the Tremendous API for creation.
// It requires the Payment field to be non-nil in CreateOrderArgs.
func (s *OrdersService) Create(args OrderArgs) (*Order, error) {
	url := fmt.Sprintf("%s/orders", s.client.BaseURL)

	// Basic validation: ensure payment is present
	if args.Payment.FundingSourceID == "" {
		return nil, errors.New("the 'payment.funding_source_id' field is required but was empty")
	}

	// Marshal the typed struct into JSON
	bodyBytes, err := json.Marshal(&args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CreateOrderArgs: %w", err)
	}

	log.Println("Request", string(bodyBytes))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header = s.client.headers()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Expect 201 Created (some APIs return 200 instead).
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var singleResp singleOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&singleResp); err != nil {
		return nil, fmt.Errorf("failed to decode create order response: %w", err)
	}
	return &singleResp.Order, nil
}
