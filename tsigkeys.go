package powerdns

import (
	"context"
	"fmt"
	"net/http"
)

// TSIGKeyService handles communication with the tsigs related methods of the Client API
type TSIGKeyService service

// TSIGKey structure with JSON API metadata
type TSIGKey struct {
	Name      *string `json:"name,omitempty"`
	ID        *string `json:"id,omitempty"`
	Algorithm *string `json:"algorithm,omitempty"`
	Key       *string `json:"key,omitempty"`
	Type      *string `json:"type,omitempty"`
}

// List retrieves a list of TSIGKeys
func (t *TSIGKeyService) List(ctx context.Context) ([]TSIGKey, error) {
	req, err := t.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("servers/%s/tsigkeys", t.client.VHost), nil, nil)
	if err != nil {
		return nil, err
	}

	tsigkeys := make([]TSIGKey, 0)
	_, err = t.client.do(req, &tsigkeys)
	return tsigkeys, err
}

// Get returns a certain TSIGKeys
func (t *TSIGKeyService) Get(ctx context.Context, id string) (*TSIGKey, error) {
	req, err := t.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("servers/%s/tsigkeys/%s", t.client.VHost, id), nil, nil)
	if err != nil {
		return nil, err
	}

	tsigkey := TSIGKey{}
	_, err = t.client.do(req, &tsigkey)
	return &tsigkey, err
}

// Create a new TSIG Key setting empty string for key will generate it
func (t *TSIGKeyService) Create(ctx context.Context, name, algorithm, key string) (*TSIGKey, error) {
	reqTsigkey := TSIGKey{
		Name:      &name,
		Algorithm: &algorithm,
		Key:       &key,
	}

	req, err := t.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("servers/%s/tsigkeys", t.client.VHost), nil, reqTsigkey)
	if err != nil {
		return nil, err
	}
	respTsigkey := TSIGKey{}
	_, err = t.client.do(req, &respTsigkey)

	return &respTsigkey, err
}

func (t *TSIGKeyService) Change(ctx context.Context, id string, newKey TSIGKey) (*TSIGKey, error) {
	req, err := t.client.newRequest(ctx, http.MethodPut, fmt.Sprintf("servers/%s/tsigkeys/%s", t.client.VHost, id), nil, newKey)
	if err != nil {
		return nil, err
	}
	responseKey := TSIGKey{}
	_, err = t.client.do(req, &responseKey)

	return &responseKey, err
}

func (t *TSIGKeyService) Delete(ctx context.Context, id string) error {
	req, err := t.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("servers/%s/tsigkeys/%s", t.client.VHost, id), nil, nil)
	if err != nil {
		return err
	}

	_, err = t.client.do(req, nil)
	return err
}
