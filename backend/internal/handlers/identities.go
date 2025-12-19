package handlers

import (
	"net/http"
	"strconv"

	"github.com/benoit-sauvere/kratos-admin-ui/backend/internal/kratos"
	"github.com/gin-gonic/gin"
	ory "github.com/ory/kratos-client-go"
)

// IdentitiesHandler handles identity-related requests
type IdentitiesHandler struct {
	client *kratos.Client
}

// NewIdentitiesHandler creates a new identities handler
func NewIdentitiesHandler(client *kratos.Client) *IdentitiesHandler {
	return &IdentitiesHandler{client: client}
}

// List returns a paginated list of identities
func (h *IdentitiesHandler) List(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	perPage, _ := strconv.ParseInt(c.DefaultQuery("per_page", "20"), 10, 64)

	result, err := h.client.ListIdentities(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch identities", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     result.Identities,
		"page":     page,
		"per_page": perPage,
		"total":    result.Total,
	})
}

// Get returns a single identity by ID
func (h *IdentitiesHandler) Get(c *gin.Context) {
	id := c.Param("id")

	identity, err := h.client.GetIdentity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Identity not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, identity)
}

// CreateIdentityRequest represents the request body for creating an identity
type CreateIdentityRequest struct {
	SchemaID string                 `json:"schema_id" binding:"required"`
	Traits   map[string]interface{} `json:"traits" binding:"required"`
	State    string                 `json:"state,omitempty"`
}

// Create creates a new identity
func (h *IdentitiesHandler) Create(c *gin.Context) {
	var req CreateIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	body := ory.CreateIdentityBody{
		SchemaId: req.SchemaID,
		Traits:   req.Traits,
	}

	if req.State != "" {
		state := ory.IdentityState(req.State)
		body.State = &state
	}

	identity, err := h.client.CreateIdentity(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create identity", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, identity)
}

// UpdateIdentityRequest represents the request body for updating an identity
type UpdateIdentityRequest struct {
	SchemaID string                 `json:"schema_id" binding:"required"`
	Traits   map[string]interface{} `json:"traits" binding:"required"`
	State    string                 `json:"state,omitempty"`
}

// Update updates an existing identity
func (h *IdentitiesHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	body := ory.UpdateIdentityBody{
		SchemaId: req.SchemaID,
		Traits:   req.Traits,
	}

	if req.State != "" {
		state := ory.IdentityState(req.State)
		body.State = state
	}

	identity, err := h.client.UpdateIdentity(c.Request.Context(), id, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update identity", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, identity)
}

// Delete deletes an identity
func (h *IdentitiesHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.client.DeleteIdentity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete identity", "details": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetSessions returns sessions for an identity
func (h *IdentitiesHandler) GetSessions(c *gin.Context) {
	id := c.Param("id")

	sessions, err := h.client.GetIdentitySessions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sessions})
}

// ResetPasswordRequest represents the request body for resetting a password
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

// ResetPassword resets the password for an identity
func (h *IdentitiesHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := h.client.ResetPassword(c.Request.Context(), id, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// DeleteCredential deletes a specific credential type for an identity
func (h *IdentitiesHandler) DeleteCredential(c *gin.Context) {
	id := c.Param("id")
	credType := c.Param("type")

	// Validate credential type
	validTypes := map[string]bool{
		"totp":          true,
		"webauthn":      true,
		"lookup_secret": true,
	}

	if !validTypes[credType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credential type", "details": "Supported types: totp, webauthn, lookup_secret"})
		return
	}

	if err := h.client.DeleteCredential(c.Request.Context(), id, credType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete credential", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credential deleted successfully"})
}

// GetWithCredentials returns a single identity by ID including credentials metadata
func (h *IdentitiesHandler) GetWithCredentials(c *gin.Context) {
	id := c.Param("id")

	identity, err := h.client.GetIdentityWithCredentials(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Identity not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, identity)
}
