package endpoints

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
)

func checkTheBody(body types.ETLPostBody) (string, error) {
	// This conditional is because we're using CreateEtl for both
	// strings (unstructured data) and JSON (structured data). The strings
	// are for SQL and starlark, and the JSON is for sequences.
	// This could be changed in time, so all ETLs are JSON, and include both
	// a "body" element and an "expected-params" element, so that we could
	// "automate" the detection of what parameters are required for each ETL.
	theBody := ""
	switch body.Body.(type) {
	case string:
		if body.Kind == "sql" || body.Kind == "starlark" {
			theBody = body.Body.(string)
		} else {
			return "", errors.New("only SQL ETL actions may come in as a string")
		}
	case []any:
		if body.Kind == "sequence" {
			// Make sure it is an array of objects with "name"
			// FIXME: the typing on these could lead to 500s if it is wrong.
			for _, e := range body.Body.([]any) {
				isOk := true
				for _, k := range []string{"name"} {
					if _, ok := e.(map[string]any)[k]; !ok {
						isOk = false
					}
				}
				if !isOk {
					return "", errors.New("missing `name` from sequence object")
				}
			}
			jsonData, err := json.MarshalIndent(body.Body, "", "  ")
			if err != nil {
				return "", errors.New("body contains invalid JSON")
			}
			theBody = string(jsonData)
		} else {
			return "", errors.New("sequences must come in as an array of JSON objects")
		}
	default:
		log.Println(body.Body)
		return "", fmt.Errorf("body contains something of an unrecognizable type: %T", body.Body)
	}

	return theBody, nil
}

func CreateEtl(c *gin.Context) {

	var body types.ETLPostBody
	// Call ShouldBindJSON to bind the incoming JSON to the newItem struct
	if err := c.ShouldBindJSON(&body); err != nil {
		// If an error occurs (e.g., invalid JSON, missing required fields),
		// return a 400 Bad Request error
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("should bind: %s", err.Error())})
		return
	}

	appId := base.GetOrPanic(c, "AppId")
	keyId := base.GetOrPanic(c, "KeyId")

	s := types.NewStorage(appId)
	err := s.SetDateYMD(body.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not parse date; must be YYYY-MM-DD, given " + body.Date,
			"date":    body.Date,
			"name":    body.Name,
		})
		return
	}

	err = s.Init()
	if err != nil {
		log.Println("storage init error: " + err.Error())
		panic(err)
	}
	// We cache whether this is loaded, so it is safe/fast to check every time
	// we try and load another ETL into the table.
	pc, _, _, _ := runtime.Caller(0)
	funcName := runtime.FuncForPC(pc).Name()
	base.LoadDefaultEtlFiles(s, funcName)

	theBody, err := checkTheBody(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": err.Error(),
			"date":    body.Date,
			"name":    body.Name,
		})
		return
	}

	if err := s.Queries.InsertETL(context.Background(), models.InsertETLParams{
		KeyID: keyId,
		Name:  body.Name,
		Kind:  body.Kind,
		Body:  sql.NullString{String: string(theBody), Valid: true},
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not save ETL step",
			"date":    body.Date,
			"name":    body.Name,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"method": c.Request.Method,
		"date":   body.Date,
		"name":   body.Name,
	})
}
