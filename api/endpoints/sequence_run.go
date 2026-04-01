package endpoints

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/types"
)

// Run a sequence immediately, synchronously.
func RunSequence(c *gin.Context) {
	appId := base.GetOrPanic(c, "AppId")
	keyId := base.GetOrPanic(c, "KeyId")
	sequenceName := c.Param("name")
	sequenceDate := c.Param("date")

	s := types.NewStorage(appId)
	s.Init()
	err := s.SetDate(sequenceDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not parse date; must be YYYY-MM-DD",
			"date":    sequenceDate,
			"name":    sequenceName,
		})
		return
	}
	errStep := ""
	seq, err := s.Queries.GetETL(context.Background(), sequenceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not find sequence",
			"date":    s.YYYYMMDD(),
			"name":    sequenceName,
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
		err := runEtl(&types.RunEtlParams{
			AppId:   appId,
			KeyId:   keyId,
			GinCtx:  nil,
			Storage: s,
			EtlName: step,
		})

		if err != nil {
			errStep = step
			log.Println("breaking on step: " + step)
			log.Println("err: " + err.Error())
			break
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": fmt.Sprintf("error in sequence execution: %s", errStep),
			"detail":  err.Error(),
			"date":    sequenceDate,
			"name":    errStep,
		})
		return
	}
}
