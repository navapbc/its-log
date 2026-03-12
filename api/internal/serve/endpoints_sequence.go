package serve

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/fsdb/models"
	"github.com/jadudm/its-log/internal/itslog"
)

func addSequenceEndpoints(rG *gin.RouterGroup) {
	auth_adminV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Admin, itslog.Test}
	auth_adminV1.Use(AuthMiddleWare(permissions))

	// Add a sequence
	auth_adminV1.POST("sequence", Sequence)
	// Run a sequence
	auth_adminV1.GET("sequence/:date/:name", Sequence)

}

// ------------------------------------------------------------------------
// Sequence
// ------------------------------------------------------------------------

type SequencePostBody struct {
	Name  string   `json:"name" binding:"required"`
	Date  string   `json:"date" binding:"required"`
	Steps []string `json:"steps" binding:"required"`
}

func Sequence(c *gin.Context) {
	// Bundle up params and call the correct method
	sctx, err := NewServeCtx(c)
	if err != nil {
		return
	}
	defer sctx.Close()

	switch c.Request.Method {
	case http.MethodPost:
		seq_post(sctx)
		return
	case http.MethodGet:
		seq_get(sctx)
		return
	case http.MethodDelete:
		// Bundle up params and call the correct method
		// seq_delete(c)
		return
	default:
		// It might not be possible to get here; Gin seems to
		// intercept unknown/underfined methods and return a 404.
		c.JSON(http.StatusBadRequest, gin.H{
			"method":  c.Request.Method,
			"message": "method not supported",
		})
		return
	}
}

func seq_post(sctx *ServeCtx) {
	var body SequencePostBody
	// Call ShouldBindJSON to bind the incoming JSON to the newItem struct
	if err := sctx.GinCtx.ShouldBindJSON(&body); err != nil {
		// If an error occurs (e.g., invalid JSON, missing required fields),
		// return a 400 Bad Request error
		sctx.GinCtx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("should bind: %s", err.Error())})
		return
	}

	// If things are malformed, return errors
	_, err := time.Parse(time.DateOnly, body.Date)
	if err != nil {
		sctx.GinCtx.JSON(http.StatusBadRequest, gin.H{
			"method":  sctx.RequestMethod,
			"message": fmt.Sprintf("%s is not YYYY-MM-DD", body.Date),
		})
		return
	}

	if err := sctx.Storage.GetQueries().InsertSequence(context.Background(), models.InsertSequenceParams{
		KeyID: sctx.KeyId,
		Name:  body.Name,
		Steps: strings.Join(body.Steps, ","),
	}); err != nil {
		sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  sctx.RequestMethod,
			"message": "could not save ETL step",
			"date":    body.Date,
			"name":    body.Name,
		})
		return
	}

	sctx.GinCtx.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"method": sctx.RequestMethod,
		"date":   body.Date,
		"name":   body.Name,
	})

}

// Run a sequence immediately, synchronously.
func seq_get(sctx *ServeCtx) {
	errStep := ""
	seq, err := sctx.Storage.GetQueries().GetSequence(context.Background(), sctx.Name)
	if err != nil {
		sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  sctx.RequestMethod,
			"message": "could not find sequence",
			"date":    sctx.YYYYMMDD(),
			"name":    sctx.Name,
		})
		return
	}

	steps := strings.Split(seq, ",")

	for _, step := range steps {
		// If anything fails, we exit the sequence
		// Pass nil so we don't write to the context in the loop
		sctx.GinCtx.Set("isSequence", true)
		err = put(sctx)
		if err != nil {
			errStep = step
			break
		}
	}
	if err != nil {
		sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  sctx.RequestMethod,
			"message": fmt.Sprintf("error in sequence execution: %s", errStep),
			"detail":  err.Error(),
			"date":    sctx.YYYYMMDD(),
			"name":    sctx.Name,
		})
		return
	}
}
