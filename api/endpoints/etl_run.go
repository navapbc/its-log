package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/types"
)

// This is the handler if we call it directly from the API.
// We only get the gin context.
func RunEtl(c *gin.Context) {
	appId := base.GetOrPanic(c, "AppId")
	keyId := base.GetOrPanic(c, "KeyId")
	date := c.Param("date")
	name := c.Param("name")

	// ETL steps can have an arbitrary POST body
	// But, they might send nothing.
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
	dateErr := s.SetDateYMD(date)
	etlP := &types.RunEtlParams{
		AppId:   appId,
		KeyId:   keyId,
		GinCtx:  c,
		Storage: s,
		EtlName: name,
		Payload: payload,
	}

	if dateErr != nil {
		ser := types.NewStandardErrorResponse(etlP, dateErr)
		ser.SetStatus(http.StatusInternalServerError).Send("could not parse date; must be YYYY-MM-DD")
		return
	}
	s.Init()

	base.LoadDefaultEtlFiles(s)

	etlErr := runEtl(etlP)

	if etlErr == nil {
		types.NewStandardOkResponse(c, etlP).SetStatus(http.StatusOK).Send()
	} else {
		ser := types.NewStandardErrorResponse(etlP, etlErr)
		ser.SetStatus(http.StatusInternalServerError).Send("golang ETL error")
	}
}

func runEtl(etlP *types.RunEtlParams) error {
	row, err := etlP.Storage.Queries.GetETL(context.Background(), etlP.EtlName)

	if err != nil {
		msg := "could not find ETL step"
		log.Println(msg)
		return fmt.Errorf("%s: %s", msg, err.Error())
	}

	switch row.Kind {
	case "sql":
		err := etlRunSql(etlP, row, nil)
		if err != nil {
			return err
		}
	case "golang":
		err := etlRunGolang(etlP, row, nil)
		if err != nil {
			return err
		}
	case "starlark":
		err := etlRunStarlark(etlP, row, nil)
		if err != nil {
			return err
		}

	case "sequence":
		// The name of the step will have been loaded into the context.

	default:
		msg := "unknown kind of etl: " + row.Kind
		log.Println(msg)
		if etlP.GinCtx != nil {
			ser := types.NewStandardErrorResponse(etlP, nil)
			ser.SetStatus(http.StatusInternalServerError).Send(msg)
		}
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	}

	updateErr := etlP.Storage.Queries.UpdateLastRun(context.Background(), etlP.EtlName)

	if updateErr != nil {
		msg := "could not update ETL metadata"
		log.Println(msg)
		if etlP.GinCtx != nil {
			ser := types.NewStandardErrorResponse(etlP, nil)
			ser.SetStatus(http.StatusInternalServerError).Send(msg)
		}
		return fmt.Errorf("%s: %s", msg, err.Error())
	}

	return nil
}
