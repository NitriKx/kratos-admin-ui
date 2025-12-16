package handlers

import (
	"net/http"
	"strconv"

	"github.com/benoit-sauvere/kratos-admin-ui/backend/internal/kratos"
	"github.com/gin-gonic/gin"
)

// SessionsHandler handles session-related requests
type SessionsHandler struct {
	client *kratos.Client
}

// NewSessionsHandler creates a new sessions handler
func NewSessionsHandler(client *kratos.Client) *SessionsHandler {
	return &SessionsHandler{client: client}
}

// List returns a paginated list of sessions
func (h *SessionsHandler) List(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	perPage, _ := strconv.ParseInt(c.DefaultQuery("per_page", "20"), 10, 64)

	sessions, err := h.client.ListSessions(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     sessions,
		"page":     page,
		"per_page": perPage,
	})
}

// Revoke revokes a session by ID
func (h *SessionsHandler) Revoke(c *gin.Context) {
	id := c.Param("id")

	if err := h.client.RevokeSession(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke session", "details": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}




