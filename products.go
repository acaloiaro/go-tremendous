package tremendous

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ProductsService handles the “products” resource.
type ProductsService struct {
	client *Client
}

// Product is a fully parsed product object (fields may vary by API).
type Product struct {
	ID          string  `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Brand       string  `json:"brand,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Currency    string  `json:"currency,omitempty"`
	Object      string  `json:"object,omitempty"`
}

// CreateProductArgs is the typed request for creating a new product.
type CreateProductArgs struct {
	Product CreateProduct `json:"product"`
}

// CreateProduct is the structure inside the top-level "product" key for creation.
type CreateProduct struct {
	Name        string  `json:"name,omitempty"`
	Brand       string  `json:"brand,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Currency    string  `json:"currency,omitempty"`
}

// singleProductResponse models the JSON for a single product spit back by the API.
type singleProductResponse struct {
	Product Product `json:"product"`
}

// listProductsResponse models the JSON structure for multiple products.
type listProductsResponse struct {
	Products []Product `json:"products"`
}

// ListProductsOptions specifies optional query filters, such as country.
type ListProductsOptions struct {
	Country string // e.g. "US" or "CA"
}

// List retrieves all products, optionally filtering by country (and other options).
func (ps *ProductsService) List(opts *ListProductsOptions) ([]Product, error) {
	base := fmt.Sprintf("%s/products", ps.client.BaseURL)

	// Build query parameters only if opts is provided
	if opts != nil {
		queryVals := url.Values{}
		if opts.Country != "" {
			queryVals.Set("country", opts.Country)
		}
		// Add additional query parameters as needed
		encoded := queryVals.Encode()
		if encoded != "" {
			base += "?" + encoded
		}
	}

	req, err := http.NewRequest(http.MethodGet, base, nil)
	if err != nil {
		return nil, err
	}
	req.Header = ps.client.headers()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var listResp listProductsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode products list: %w", err)
	}
	return listResp.Products, nil
}

// Create adds a new product to the Tremendous API.
func (ps *ProductsService) Create(args CreateProductArgs) (*Product, error) {
	url := fmt.Sprintf("%s/products", ps.client.BaseURL)

	jsonBytes, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CreateProductArgs: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header = ps.client.headers()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Expect 201 Created or 200 is possible in some APIs
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var singleResp singleProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&singleResp); err != nil {
		return nil, fmt.Errorf("failed to decode create product response: %w", err)
	}
	return &singleResp.Product, nil
}

// Retrieve fetches a single product by its ID.
func (ps *ProductsService) Retrieve(productID string) (*Product, error) {
	productID = strings.TrimSpace(productID)
	url := fmt.Sprintf("%s/products/%s", ps.client.BaseURL, productID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = ps.client.headers()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var singleResp singleProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&singleResp); err != nil {
		return nil, fmt.Errorf("failed to decode single product: %w", err)
	}
	return &singleResp.Product, nil
}
