package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const defaultBaseURL = "https://www.dati.gov.it/opendata/api/3/action"

// Client è l'HTTP client per l'API CKAN di dati.gov.it.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New crea un Client con il base URL di dati.gov.it e timeout 30s.
func New() *Client {
	return NewWithBase(defaultBaseURL)
}

// NewWithBase crea un Client con un base URL personalizzato (utile per i test).
func NewWithBase(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SearchDatasets chiama package_search.
func (c *Client) SearchDatasets(ctx context.Context, query, tags, org string, rows int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("rows", strconv.Itoa(rows))
	if tags != "" {
		params.Set("fq", fmt.Sprintf("tags:%s", tags))
	}
	if org != "" {
		params.Set("fq", fmt.Sprintf("organization:%s", org))
	}
	return doRequest[SearchResult](ctx, c, "package_search", params)
}

// GetDataset chiama package_show.
func (c *Client) GetDataset(ctx context.Context, id string) (*Dataset, error) {
	return doRequest[Dataset](ctx, c, "package_show", url.Values{"id": {id}})
}

// ListOrganizations chiama organization_list.
func (c *Client) ListOrganizations(ctx context.Context) ([]string, error) {
	result, err := doRequest[[]string](ctx, c, "organization_list", url.Values{})
	if err != nil {
		return nil, err
	}
	return *result, nil
}

// GetOrganization chiama organization_show.
func (c *Client) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	return doRequest[Organization](ctx, c, "organization_show", url.Values{"id": {id}})
}

// doRequest esegue la chiamata HTTP con retry 3x e backoff esponenziale.
func doRequest[T any](ctx context.Context, c *Client, action string, params url.Values) (*T, error) {
	endpoint := fmt.Sprintf("%s/%s?%s", c.baseURL, action, params.Encode())
	var lastErr error

	for attempt := range 3 {
		if attempt > 0 {
			wait := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("rate limited (HTTP 429)")
			continue
		}

		var envelope CKANResponse[T]
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		if !envelope.Success {
			if envelope.Error != nil {
				return nil, fmt.Errorf("CKAN error: %s", envelope.Error.Message)
			}
			return nil, fmt.Errorf("CKAN request failed (success=false)")
		}
		return &envelope.Result, nil
	}
	return nil, fmt.Errorf("dopo 3 tentativi: %w", lastErr)
}
