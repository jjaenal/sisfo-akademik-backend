package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/service"
)

type schoolServiceImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewSchoolService(baseURL string, timeout time.Duration) service.SchoolService {
	return &schoolServiceImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *schoolServiceImpl) GetLocation(ctx context.Context, tenantID string) (*service.SchoolLocation, error) {
	url := fmt.Sprintf("%s/api/v1/schools/tenant/%s", s.baseURL, tenantID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			// We can assume a default radius or add it to the School entity in academic-service later.
			// For now, we'll return default radius if not present in response, 
			// but since the SchoolLocation struct has Radius, we should try to map it if available.
			// The academic-service School entity currently doesn't seem to have Radius.
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("api returned failure")
	}

	return &service.SchoolLocation{
		Latitude:  response.Data.Latitude,
		Longitude: response.Data.Longitude,
		Radius:    100.0, // Default radius 100 meters, as it is not yet in academic-service
	}, nil
}
