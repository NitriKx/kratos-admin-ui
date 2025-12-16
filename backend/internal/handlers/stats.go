package handlers

import (
	"net/http"

	"github.com/benoit-sauvere/kratos-admin-ui/backend/internal/kratos"
	"github.com/gin-gonic/gin"
)

// StatsHandler handles stats-related requests
type StatsHandler struct {
	client *kratos.Client
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(client *kratos.Client) *StatsHandler {
	return &StatsHandler{client: client}
}

// StatsResponse represents dashboard statistics
type StatsResponse struct {
	ActiveIdentities int64 `json:"active_identities"`
	ActiveSessions   int64 `json:"active_sessions"`
}

// Get returns dashboard statistics
func (h *StatsHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	// Get active identity count
	activeIdentities, err := h.client.GetActiveIdentityCount(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active identity count", "details": err.Error()})
		return
	}

	// Get active session count
	activeSessions, err := h.client.GetSessionCount(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch session count", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, StatsResponse{
		ActiveIdentities: activeIdentities,
		ActiveSessions:   activeSessions,
	})
}

