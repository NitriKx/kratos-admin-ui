package handlers

import (
	"net/http"

	"github.com/benoit-sauvere/kratos-admin-ui/backend/internal/kratos"
	"github.com/gin-gonic/gin"
)

// SchemasHandler handles schema-related requests
type SchemasHandler struct {
	client *kratos.Client
}

// NewSchemasHandler creates a new schemas handler
func NewSchemasHandler(client *kratos.Client) *SchemasHandler {
	return &SchemasHandler{client: client}
}

// List returns all identity schemas
func (h *SchemasHandler) List(c *gin.Context) {
	schemas, err := h.client.ListIdentitySchemas(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schemas", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schemas})
}




