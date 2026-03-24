package serve

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/navapbc/its-log/internal/base"

	"github.com/gin-gonic/gin"
)

func addSequenceEndpoints(rG *gin.RouterGroup) {
	auth_adminV1 := rG.Group("/")
	permissions := []base.PermissionType{base.Admin, base.Test}
	auth_adminV1.Use(AuthMiddleWare(permissions))

	// Insert a sequence
	// It's actually just inserting an entry into the
	// ETL table with the correct values. We give it a 'sequence'
	// endpoint, but it is interchangeable.
	auth_adminV1.POST("sequence/:date", ETL)

	// Run a sequence
	auth_adminV1.GET("sequence/:date/:name", RunSequence)

}

// ------------------------------------------------------------------------
// Sequence
// ------------------------------------------------------------------------

// Run a sequence immediately, synchronously.
func RunSequence(c *gin.Context) {
	sctx, err := base.NewServeCtx(c)
	if err != nil {
		return
	}
	defer sctx.Close()

	errStep := ""
	seq, err := sctx.Storage.GetQueries().GetETL(context.Background(), sctx.Name)
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

	// Split on commas if they are present; otherwise
	// split on newlines.
	trimmed := strings.TrimSpace(seq.Body.String)
	var steps []string
	if strings.Contains(trimmed, ",") {
		steps = strings.Split(trimmed, ",")
	} else {
		steps = strings.Split(trimmed, "\n")
	}

	for _, step := range steps {
		// If anything fails, we exit the sequence
		// Pass nil so we don't write to the context in the loop
		sctx.GinCtx.Set("isSequence", true)
		sctx.Name = step

		// log.Println("before")
		// resp := make(chan int64)
		// base.FunsQ <- base.Job{
		// 	Id: fmt.Sprintf("RunSeq: %s", step),
		// 	Op: func() {
		// 		log.Println("running put")
		// 		err = put(sctx)
		// 	},
		// 	Resp: resp,
		// }
		// v := <-resp
		// log.Printf("after: %d\n", v)

		err = put(sctx)

		if err != nil {
			errStep = step
			log.Println("breaking on step: " + step)
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
