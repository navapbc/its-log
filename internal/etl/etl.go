package etl

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
)

type ActionHandler func(gjson.Result, *sql.DB) error
type ActionMap map[string]ActionHandler

var actionMap ActionMap = make(ActionMap, 0)

func initMap() {
	actionMap["message"] = Message
	actionMap["sql"] = Sql
	actionMap["fileCopy"] = FileCopy
}

func Run(runscript string, sql *sql.DB) {
	initMap()
	actions := gjson.Get(runscript, "actions")
	actions.ForEach(func(ndx, action gjson.Result) bool {
		// println(key.String() + " | " + value.String())
		err := actionMap[action.Get("action").String()](action, sql)
		// Continue as long as we got a nil result
		return err == nil
	})
}

func Message(action gjson.Result, sql *sql.DB) error {
	fmt.Printf("== message ==\n%s\n", action.Get("message").String())
	return nil
}

func Sql(action gjson.Result, sql *sql.DB) error {
	fmt.Printf("== sql: %s ==\n", action.Get("filename").String())
	query, err := os.ReadFile(action.Get("filename").String())
	if err != nil {
		panic(err)
	}
	_, err = sql.Exec(string(query))
	if err != nil {
		panic(err)
	}
	return nil
}

func FileCopy(action gjson.Result, sql *sql.DB) error {
	fmt.Printf("== file copy ==\n  src: %s\n  dst: %s\n", action.Get("source").String(), action.Get("destination").String())
	return nil
}
