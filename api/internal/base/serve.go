package base

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

const OK = "ok"
const ERROR = "error"

type Success struct {
	Status string `json:"status"`
}
type Error struct {
	Status    string         `json:"status"`
	Data      map[string]any `json:"data"`
	ErrorType string         `json:"error_type"`
	Error     string         `json:"error"`
}

func GetOrPanic(c *gin.Context, key string) string {
	v, exists := c.Get(key)
	if !exists {
		panic(fmt.Sprintf("could not get key from gin context: %s", key))
	}
	return v.(string)
}

// -----------------------------
// ServeCtx
// Every call has context it needs to execute.
// This bundles it up in one place to eliminate
// copypasta in the endpoint handling.
type ServeCtx struct {
	AppId         string
	KeyId         string
	Storage       *SqliteStorage
	RequestMethod string
	Date          time.Time
	Name          string
	GinCtx        *gin.Context
}

func NewServeCtx(c *gin.Context) (*ServeCtx, error) {
	etlctx := &ServeCtx{}
	err := etlctx.Init(c)
	if err != nil {
		return nil, err
	}
	return etlctx, nil
}

func (ectx *ServeCtx) Init(c *gin.Context) error {
	ectx.GinCtx = c

	// Extract the app and key id from the middleware context
	ectx.AppId = GetOrPanic(c, ITSLOG_APPID)
	ectx.KeyId = GetOrPanic(c, ITSLOG_KEYID)
	// Extract the date/name from the URL
	ectx.Name = c.Param("name")
	param_date := c.Param("date")

	if param_date != "" {
		// Convert the date to a time.Time structure
		parsed_date, err := time.Parse(time.DateOnly, param_date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"method":  c.Request.Method,
				"message": fmt.Sprintf("%s is not YYYY-MM-DD", param_date),
			})
			return err
		}
		ectx.Date = parsed_date
	}

	// Setup the storage
	storage := &SqliteStorage{
		AppId: ectx.AppId,
		Date:  ectx.YYYYMMDD(),
		Kind:  NamedDatabase,
		Path:  viper.GetString("storage.path"),
	}

	err := storage.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"method":  c.Request.Method,
			"message": fmt.Sprintf("could not initialize database; %s-%s", ectx.YYYYMMDD(), ectx.AppId),
		})
		return err
	}
	ectx.Storage = storage

	// Pull some things from the context for use/reference
	ectx.RequestMethod = c.Request.Method
	return nil
}

func (ectx *ServeCtx) YYYYMMDD() string {
	return fmt.Sprintf("%s", ectx.Date.Format(time.DateOnly))
}

func (ectx *ServeCtx) Close() {
	ectx.Storage.Close()
}
