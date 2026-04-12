package types

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StandardErrorResponseParams struct {
	c      *gin.Context
	err    error
	etlP   *RunEtlParams
	Status int
	Method string
	Date   ILTime
	Name   string
}

func NewStandardErrorResponse(etlP *RunEtlParams, err error) *StandardErrorResponseParams {
	return &StandardErrorResponseParams{
		c:    etlP.GinCtx,
		err:  err,
		etlP: etlP,
	}
}

func (ser *StandardErrorResponseParams) SetStatus(status int) *StandardErrorResponseParams {
	ser.Status = status
	return ser
}

func (ser *StandardErrorResponseParams) Send(msg string) {

	zap.L().Error("ERR",
		zap.Int("status", ser.Status),
		zap.String("msg", msg),
		zap.String("error", ser.err.Error()),
		zap.String("name", ser.etlP.EtlName),
		zap.String("date", ser.etlP.Storage.YYYYMMDD()),
	)
	ser.c.JSON(ser.Status, gin.H{
		"status":  "error",
		"method":  ser.c.Request.Method,
		"message": fmt.Sprintf("%s: %s", msg, ser.err.Error()),
		"date":    ser.etlP.Storage.YYYYMMDD(),
		"name":    ser.etlP.EtlName,
	})
	return
}

type StandardOkResponseParams struct {
	c      *gin.Context
	etlP   *RunEtlParams
	Status int
}

func NewStandardOkResponse(c *gin.Context, etlP *RunEtlParams) *StandardOkResponseParams {
	return &StandardOkResponseParams{
		c:    c,
		etlP: etlP,
	}
}

func (sok *StandardOkResponseParams) SetStatus(status int) *StandardOkResponseParams {
	sok.Status = status
	return sok
}

func (sok *StandardOkResponseParams) Send() {
	name := sok.etlP.EtlName
	date := sok.etlP.Storage.ILTime.AsYYYYMMDD()

	zap.L().Info("OK",
		zap.Int("status", sok.Status),
		zap.String("name", name),
		zap.String("date", date),
	)

	sok.c.JSON(sok.Status, gin.H{
		"status": "ok",
		"method": sok.c.Request.Method,
		"date":   date,
		"name":   name,
	})
}
