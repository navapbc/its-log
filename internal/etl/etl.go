package etl

import (
	"context"
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
	actionMap["assert"] = Assert
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

	// Get a Tx for making transaction requests.
	ctx := context.Background()
	tx, err := sql.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, string(query))
	if err != nil {
		panic(err)
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		panic(err)
	}
	return nil
}

func Assert(action gjson.Result, sql *sql.DB) error {
	fmt.Printf("== assert: %s ==\n", action.Get("filename").String())
	query, err := os.ReadFile(action.Get("filename").String())
	if err != nil {
		panic(err)
	}

	// Get a Tx for making transaction requests.
	ctx := context.Background()
	tx, err := sql.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	res := tx.QueryRowContext(ctx, string(query))
	var found bool
	err = res.Scan(&found)
	if err != nil {
		panic(err)
	}
	if !found {
		fmt.Println("assertion returned false")
		tx.Rollback()
		os.Exit(-1)
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		panic(err)
	}
	return nil
}

func FileCopy(action gjson.Result, sql *sql.DB) error {
	fmt.Printf("== file copy ==\n  src: %s\n  dst: %s\n", action.Get("source").String(), action.Get("destination").String())
	return nil
}
