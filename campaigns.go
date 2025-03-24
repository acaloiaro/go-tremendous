package tremendous

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CampaignsService handles the "campaigns" resource.
type CampaignsService struct {
	client *Client
}

// Campaign represents a single campaign record from Tremendous.
type Campaign struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

// listCampaignsResponse models the JSON response for GET /campaigns.
type listCampaignsResponse struct {
	Campaigns []Campaign `json:"campaigns"`
}

// List fetches a filtered list of campaigns from the API.
//
// Example usage:
//
//	client.Campaigns.List(&tremendous.ListCampaignsOptions{Name: "Summer", Status: "active"})
func (cs *CampaignsService) List() ([]Campaign, error) {
	// Construct the base endpoint
	baseURL := fmt.Sprintf("%s/campaigns", cs.client.BaseURL)

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = cs.client.headers()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var answer listCampaignsResponse
	if err := json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		return nil, fmt.Errorf("failed to decode campaigns list: %w", err)
	}
	return answer.Campaigns, nil
}

// Retrieve fetches a single campaign by its ID if needed.
// (Optional: Implement similarly to the templates/products retrieve.)
func (cs *CampaignsService) Retrieve(campaignID string) (*Campaign, error) {
	campaignID = strings.TrimSpace(campaignID)
	urlStr := fmt.Sprintf("%s/campaigns/%s", cs.client.BaseURL, campaignID)

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header = cs.client.headers()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var single struct {
		Campaign Campaign `json:"campaign"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&single); err != nil {
		return nil, fmt.Errorf("failed to decode single campaign: %w", err)
	}
	return &single.Campaign, nil
}
