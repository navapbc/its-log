package serve

import (
	_ "github.com/creasty/defaults"
)

// // -----------------------------
// // ServeCtx
// // Every call has context it needs to execute.
// // This bundles it up in one place to eliminate
// // copypasta in the endpoint handling.
// type ServeCtx struct {
// 	AppId         string
// 	KeyId         string
// 	Storage       *fsdb.SqliteStorage
// 	RequestMethod string
// 	Date          time.Time
// 	Name          string
// 	GinCtx        *gin.Context
// }

// func NewServeCtx(c *gin.Context) (*ServeCtx, error) {
// 	etlctx := &ServeCtx{}
// 	err := etlctx.Init(c)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return etlctx, nil
// }

// func (ectx *ServeCtx) Init(c *gin.Context) error {
// 	ectx.GinCtx = c

// 	// Extract the app and key id from the middleware context
// 	ectx.AppId = types.GetOrPanic(c, types.ITSLOG_APPID)
// 	ectx.KeyId = types.GetOrPanic(c, types.ITSLOG_KEYID)
// 	// Extract the date/name from the URL
// 	ectx.Name = c.Param("name")
// 	param_date := c.Param("date")

// 	if param_date != "" {
// 		// Convert the date to a time.Time structure
// 		parsed_date, err := time.Parse(time.DateOnly, param_date)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"method":  c.Request.Method,
// 				"message": fmt.Sprintf("%s is not YYYY-MM-DD", param_date),
// 			})
// 			return err
// 		}
// 		ectx.Date = parsed_date
// 	}

// 	// Setup the storage
// 	storage := &fsdb.SqliteStorage{
// 		AppId: ectx.AppId,
// 		Date:  ectx.YYYYMMDD(),
// 		Kind:  fsdb.NamedDatabase,
// 		Path:  viper.GetString("storage.path"),
// 	}

// 	err := storage.Init()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"method":  c.Request.Method,
// 			"message": fmt.Sprintf("could not initialize database; %s-%s", ectx.YYYYMMDD(), ectx.AppId),
// 		})
// 		return err
// 	}
// 	ectx.Storage = storage

// 	// Pull some things from the context for use/reference
// 	ectx.RequestMethod = c.Request.Method
// 	return nil
// }

// func (ectx *ServeCtx) YYYYMMDD() string {
// 	return fmt.Sprintf("%s", ectx.Date.Format(time.DateOnly))
// }

// func (ectx *ServeCtx) Close() {
// 	ectx.Storage.Close()
// }
