package etl

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type EtlParams struct {
	EtlRunscriptPath string
	EtlDate          string
	EtlApiKey        string
	EtlUrl           string
}

type ActionHandler func(gjson.Result, EtlParams) error
type ActionMap map[string]ActionHandler

var actionMap ActionMap = make(ActionMap, 0)

func initMap() {
	actionMap["message"] = Message
	actionMap["run"] = RunSql
	actionMap["load"] = LoadSql
	actionMap["assert"] = Assert
}

func Run(runscript string, params EtlParams) {
	initMap()
	actions := gjson.Get(runscript, "actions")
	actions.ForEach(func(ndx, action gjson.Result) bool {
		// println(key.String() + " | " + value.String())
		err := actionMap[action.Get("action").String()](action, params)
		// Continue as long as we got a nil result
		return err == nil
	})
}

func Message(action gjson.Result, params EtlParams) error {
	fmt.Printf("== message ==\n%s\n", action.Get("message").String())
	return nil
}

func RunSql(action gjson.Result, params EtlParams) error {

	// FIXME: This is only for local/debugging.
	// In production, we would not want to do this. Or, it should be
	// a command-line override
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	name := action.Get("name").String()
	url := fmt.Sprintf("%s/%s/%s", params.EtlUrl, params.EtlDate, name)
	log.Println(url)

	// bytes.NewBuffer(payload)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		log.Fatal("Failed to create request object: " + err.Error())
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", params.EtlApiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to execute client: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read body: " + err.Error())
	}

	// FIXME better handling here; this is essentially printf debugging
	log.Println(string(body))

	return nil
}

func LoadSql(action gjson.Result, params EtlParams) error {
	filename := action.Get("filename").String()
	name := action.Get("name").String()
	dir := filepath.Dir(params.EtlRunscriptPath)
	with_dir := filepath.Join(dir, filename)

	contents, err := os.ReadFile(with_dir)
	if err != nil {
		log.Fatal("could not read file: " + with_dir)
	}

	jsonData, _ := sjson.Set("", "sql", contents)
	var jsonStr = []byte(jsonData)
	payload := bytes.NewBuffer(jsonStr)

	// FIXME: This is only for local/debugging.
	// In production, we would not want to do this. Or, it should be
	// a command-line override
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := fmt.Sprintf("%s/%s/%s", params.EtlUrl, params.EtlDate, name)
	log.Println(url)

	// bytes.NewBuffer(payload)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		log.Fatal("Failed to create request object: " + err.Error())
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", params.EtlApiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to execute client: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read body: " + err.Error())
	}

	// FIXME better handling here; this is essentially printf debugging
	log.Println(string(body))

	return nil
}

func Assert(action gjson.Result, params EtlParams) error {
	fmt.Printf("== assert: %s ==\n", action.Get("filename").String())
	// // query, err := os.ReadFile(action.Get("filename").String())
	// if err != nil {
	// 	panic(err)
	// }

	// // Get a Tx for making transaction requests.
	// ctx := context.Background()
	// tx, err := sql.BeginTx(ctx, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// // Defer a rollback in case anything fails.
	// defer tx.Rollback()

	// res := tx.QueryRowContext(ctx, string(query))
	// var found bool
	// err = res.Scan(&found)
	// if err != nil {
	// 	panic(err)
	// }
	// if !found {
	// 	fmt.Println("assertion returned false")
	// 	tx.Rollback()
	// 	os.Exit(-1)
	// }
	// // Commit the transaction.
	// if err = tx.Commit(); err != nil {
	// 	panic(err)
	// }

	return nil
}
