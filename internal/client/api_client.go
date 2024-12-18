package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"queue-manager/internal/models"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *APIClient) CreateTenant(clientID, name string) error {
	payload := models.CreateTenantRequest{
		ClientID: clientID,
		Name:     name,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/tenants", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create tenant: unexpected status code %d", resp.StatusCode)
	}

	return nil
}

func (c *APIClient) DeleteTenant(clientID string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/tenants/%s", c.baseURL, clientID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete tenant: unexpected status code %d", resp.StatusCode)
	}

	return nil
}

func (c *APIClient) ProcessPayload(clientID string, payload string) error {
	var jsonPayload interface{}
	if err := json.Unmarshal([]byte(payload), &jsonPayload); err != nil {
		return fmt.Errorf("invalid JSON payload: %w", err)
	}

	request := models.ProcessPayloadRequest{
		ClientID: clientID,
		Payload:  jsonPayload,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/tenants/process", c.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to process payload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to process payload: unexpected status code %d", resp.StatusCode)
	}

	return nil
} 