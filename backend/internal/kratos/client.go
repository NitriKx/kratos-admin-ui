package kratos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	ory "github.com/ory/kratos-client-go"
)

// Client wraps the Ory Kratos admin API client
type Client struct {
	api       *ory.APIClient
	publicURL string
}

// NewClient creates a new Kratos client
func NewClient(adminURL string) *Client {
	config := ory.NewConfiguration()
	config.Servers = []ory.ServerConfiguration{
		{URL: adminURL},
	}

	return &Client{
		api: ory.NewAPIClient(config),
	}
}

// SetPublicURL sets the public URL for the client
func (c *Client) SetPublicURL(publicURL string) {
	c.publicURL = publicURL
}

// ListIdentitiesResult contains the list of identities and pagination info
type ListIdentitiesResult struct {
	Identities []ory.Identity
	Total      int64
}

// ListIdentities retrieves a paginated list of identities
func (c *Client) ListIdentities(ctx context.Context, page, perPage int64) (*ListIdentitiesResult, error) {
	// First, get the total count by fetching with a large page size
	// Kratos doesn't have a native count endpoint, so we need this workaround
	allIdentities, _, err := c.api.IdentityApi.ListIdentities(ctx).PerPage(10000).Execute()
	if err != nil {
		return nil, err
	}
	total := int64(len(allIdentities))

	// Calculate offset for pagination
	offset := (page - 1) * perPage
	end := offset + perPage
	if end > total {
		end = total
	}

	// Return the slice for the requested page
	var pageIdentities []ory.Identity
	if offset < total {
		pageIdentities = allIdentities[offset:end]
	} else {
		pageIdentities = []ory.Identity{}
	}

	return &ListIdentitiesResult{
		Identities: pageIdentities,
		Total:      total,
	}, nil
}

// GetIdentity retrieves a single identity by ID
func (c *Client) GetIdentity(ctx context.Context, id string) (*ory.Identity, error) {
	identity, _, err := c.api.IdentityApi.GetIdentity(ctx, id).Execute()
	if err != nil {
		return nil, err
	}

	return identity, nil
}

// CreateIdentity creates a new identity
func (c *Client) CreateIdentity(ctx context.Context, body ory.CreateIdentityBody) (*ory.Identity, error) {
	identity, _, err := c.api.IdentityApi.CreateIdentity(ctx).CreateIdentityBody(body).Execute()
	if err != nil {
		return nil, err
	}

	return identity, nil
}

// UpdateIdentity updates an existing identity
func (c *Client) UpdateIdentity(ctx context.Context, id string, body ory.UpdateIdentityBody) (*ory.Identity, error) {
	identity, _, err := c.api.IdentityApi.UpdateIdentity(ctx, id).UpdateIdentityBody(body).Execute()
	if err != nil {
		return nil, err
	}

	return identity, nil
}

// DeleteIdentity deletes an identity
func (c *Client) DeleteIdentity(ctx context.Context, id string) error {
	_, err := c.api.IdentityApi.DeleteIdentity(ctx, id).Execute()
	return err
}

// GetIdentitySessions retrieves sessions for an identity
func (c *Client) GetIdentitySessions(ctx context.Context, id string) ([]ory.Session, error) {
	sessions, _, err := c.api.IdentityApi.ListIdentitySessions(ctx, id).Execute()
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// ListSessions retrieves all sessions
func (c *Client) ListSessions(ctx context.Context, page, perPage int64) ([]ory.Session, error) {
	req := c.api.IdentityApi.ListSessions(ctx)
	req = req.PageSize(perPage)

	sessions, _, err := req.Execute()
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// RevokeSession revokes a session by ID
func (c *Client) RevokeSession(ctx context.Context, id string) error {
	_, err := c.api.IdentityApi.DisableSession(ctx, id).Execute()
	return err
}

// IdentitySchemaWithContent represents a schema with its full content
type IdentitySchemaWithContent struct {
	ID     string                 `json:"id"`
	Schema map[string]interface{} `json:"schema"`
}

// ListIdentitySchemas retrieves all identity schemas from the public API
func (c *Client) ListIdentitySchemas(ctx context.Context) ([]IdentitySchemaWithContent, error) {
	if c.publicURL == "" {
		return nil, fmt.Errorf("public URL not configured")
	}

	// Fetch schema list from public API - it includes full schema content inline
	resp, err := http.Get(c.publicURL + "/schemas")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schemas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var schemas []IdentitySchemaWithContent
	if err := json.NewDecoder(resp.Body).Decode(&schemas); err != nil {
		return nil, fmt.Errorf("failed to decode schemas: %w", err)
	}

	return schemas, nil
}

// GetIdentityCount returns the total count of identities
func (c *Client) GetIdentityCount(ctx context.Context) (int64, error) {
	// Kratos doesn't have a direct count endpoint, so we fetch all identities
	identities, _, err := c.api.IdentityApi.ListIdentities(ctx).PerPage(10000).Execute()
	if err != nil {
		return 0, err
	}

	return int64(len(identities)), nil
}

// GetActiveIdentityCount returns the count of active identities
func (c *Client) GetActiveIdentityCount(ctx context.Context) (int64, error) {
	// Fetch all identities and count active ones
	identities, _, err := c.api.IdentityApi.ListIdentities(ctx).PerPage(10000).Execute()
	if err != nil {
		return 0, err
	}

	var count int64
	for _, identity := range identities {
		if identity.State != nil && string(*identity.State) == "active" {
			count++
		}
	}

	return count, nil
}

// GetSessionCount returns the total count of active sessions
func (c *Client) GetSessionCount(ctx context.Context) (int64, error) {
	sessions, _, err := c.api.IdentityApi.ListSessions(ctx).Active(true).Execute()
	if err != nil {
		return 0, err
	}

	return int64(len(sessions)), nil
}

// ResetPassword sets a new password for an identity
func (c *Client) ResetPassword(ctx context.Context, id string, newPassword string) error {
	// First, get the current identity to preserve its data
	identity, _, err := c.api.IdentityApi.GetIdentity(ctx, id).Execute()
	if err != nil {
		return fmt.Errorf("failed to get identity: %w", err)
	}

	// Convert traits to map[string]interface{}
	traits, ok := identity.Traits.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to convert traits to map")
	}

	// Build update body with new password credentials
	body := ory.UpdateIdentityBody{
		SchemaId: identity.SchemaId,
		Traits:   traits,
		State:    *identity.State,
		Credentials: &ory.IdentityWithCredentials{
			Password: &ory.IdentityWithCredentialsPassword{
				Config: &ory.IdentityWithCredentialsPasswordConfig{
					Password: &newPassword,
				},
			},
		},
	}

	_, _, err = c.api.IdentityApi.UpdateIdentity(ctx, id).UpdateIdentityBody(body).Execute()
	if err != nil {
		return fmt.Errorf("failed to update identity with new password: %w", err)
	}

	return nil
}

// DeleteCredential deletes a specific credential type for an identity
func (c *Client) DeleteCredential(ctx context.Context, id string, credentialType string) error {
	_, err := c.api.IdentityApi.DeleteIdentityCredentials(ctx, id, credentialType).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete %s credentials: %w", credentialType, err)
	}
	return nil
}

// GetIdentityWithCredentials retrieves a single identity by ID including credentials metadata
func (c *Client) GetIdentityWithCredentials(ctx context.Context, id string) (*ory.Identity, error) {
	identity, _, err := c.api.IdentityApi.GetIdentity(ctx, id).IncludeCredential([]string{"totp", "password", "oidc", "webauthn", "lookup_secret"}).Execute()
	if err != nil {
		return nil, err
	}

	return identity, nil
}
