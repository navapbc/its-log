package types

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type ETLPostBody struct {
	Name string `json:"name" binding:"required"`
	Kind string `json:"kind" binding:"required"`
	Body any    `json:"body" binding:"required"`
	Date string `json:"date" binding:"required"`
}

type RunEtlParams struct {
	AppId   string
	GinCtx  *gin.Context
	EtlName string
	KeyId   string
	Storage *Storage
	Payload map[string]any
}

type SummaryRow struct {
	ID        int64          `json:"id"`
	LastRun   int64          `json:"last_run"`
	Date      string         `json:"date"`
	KeyId     string         `json:"key_id"`
	Operation string         `json:"operation"`
	Tags      sql.NullString `json:"tags"`
	Value     sql.NullString `json:"value"`
	Count     float64        `json:"count"`
	Hash      sql.NullString `json:"hash"`
}
