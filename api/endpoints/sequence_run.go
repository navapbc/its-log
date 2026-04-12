package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"

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

	// // ETL steps can have an arbitrary POST body?
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		body = []byte("{}")
	}
	payload := make(map[string]any)
	jsonErr := json.Unmarshal(body, &payload)
	if jsonErr != nil {
		payload = make(map[string]any)
	}

	s := types.NewStorage(appId)
	dateErr := s.SetDateYMD(sequenceDate)
	if dateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not parse date; must be YYYY-MM-DD",
			"date":    sequenceDate,
			"name":    sequenceName,
		})
		return
	}
	s.Init()

	pc, _, _, _ := runtime.Caller(0)
	funcName := runtime.FuncForPC(pc).Name()
	base.LoadDefaultEtlFiles(s, funcName)

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

	// Sequences are JSON.
	// [ { name: "count-total" }, { name: "consolidate", params: { ... } } ]
	var seqEntries []types.SequenceEntry
	jsonErr2 := json.Unmarshal([]byte(seq.Body.String), &seqEntries)
	if jsonErr2 != nil {
		msg := fmt.Sprintf("malformed sequence spec: %s", seq.Body.String)
		log.Println(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": msg,
			"date":    s.YYYYMMDD(),
			"name":    sequenceName,
		})
		return
	}

	for _, step := range seqEntries {
		// Merge the passed-in payload with whatever is defined for this step.
		merged := make(map[string]any)
		for k, v := range payload {
			merged[k] = v
		}
		// DEBUG LOG
		// log.Printf("step params: %v\n", step.Params)
		if step.Params != nil {
			for k, v := range step.Params {
				merged[k] = v
			}
		}

		// DEBUG LOG
		// log.Printf("date[%s] sequence[%s] step[%s]\n", s.YYYYMMDD(), sequenceName, step.Name)

		sequenceError := runEtl(&types.RunEtlParams{
			AppId:   appId,
			KeyId:   keyId,
			GinCtx:  nil,
			Storage: s,
			EtlName: step.Name,
			Payload: merged,
		})

		if sequenceError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  c.Request.Method,
				"message": fmt.Sprintf("error in sequence execution: %s", step.Name),
				"detail":  sequenceError.Error(),
				"date":    sequenceDate,
				"name":    step.Name,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
