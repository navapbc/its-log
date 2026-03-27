package types

import (
	"github.com/gin-gonic/gin"
)

type ETLPostBody struct {
	Name string `json:"name" binding:"required"`
	Kind string `json:"kind" binding:"required"`
	Body string `json:"body" binding:"required"`
	Date string `json:"date" binding:"required"`
}

type RunEtlParams struct {
	AppId   string
	GinCtx  *gin.Context
	EtlName string
	KeyId   string
	Storage *Storage
}
