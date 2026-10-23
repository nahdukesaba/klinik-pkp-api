package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CentralProvisioner struct {
	baseURL string
	client  *http.Client
}

func NewCentralProvisioner(baseURL string) *CentralProvisioner {
	return &CentralProvisioner{baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"), client: &http.Client{Timeout: 15 * time.Second}}
}

func (p *CentralProvisioner) Provision(ctx context.Context, authorization string, input DirectoryCreate) (uuid.UUID, error) {
	if p == nil || p.baseURL == "" || strings.TrimSpace(authorization) == "" {
		return uuid.Nil, ErrProvisioningUnavailable
	}
	body, err := json.Marshal(input)
	if err != nil {
		return uuid.Nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/v1/admin/klinik/users", bytes.NewReader(body))
	if err != nil {
		return uuid.Nil, err
	}
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrProvisioningUnavailable, err)
	}
	defer resp.Body.Close()
	var payload struct {
		Data struct {
			AppUserID uuid.UUID `json:"appUserId"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid Central response", ErrProvisioningUnavailable)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return uuid.Nil, fmt.Errorf("%w: Central returned status %d", ErrProvisioningUnavailable, resp.StatusCode)
	}
	if payload.Data.AppUserID == uuid.Nil {
		return uuid.Nil, ErrProvisioningUnavailable
	}
	return payload.Data.AppUserID, nil
}
