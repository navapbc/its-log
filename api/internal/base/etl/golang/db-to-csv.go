package etl

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

// Write sql `rows` as a csv to the `writer`.
// table: the table within the query, for logging
func rowsToCSV(
	rows *sql.Rows,
	table string,
	writer *csv.Writer,
) error {
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("could not read columns: %s", table)
	}
	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("could not write column names: %s", table)
	}

	for rows.Next() {
		values := make([]any, len(columns))
		valuePointers := make([]any, len(columns))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return fmt.Errorf("could not read row: %s", table)
		}

		csvRow := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				csvRow[i] = "" // Handle NULL values
			} else {
				csvRow[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := writer.Write(csvRow); err != nil {
			return fmt.Errorf("could not write row: %s", table)
		}
	}

	if err = rows.Err(); err != nil {
		return errors.New("row handling error at end of process")
	}

	return nil
}

func queryToCSV(storage *types.Storage, query string, table string, path string) error {
	rows, err := storage.GetDB().Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create CSV: %s", path)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = rowsToCSV(rows, table, writer)
	if err != nil {
		return err
	}

	log.Println("exported: " + path)

	return nil
}

func SqliteToCSV(etlP *types.RunEtlParams) error {
	// var sqliteToCsvPrefix string = "sqlite_to_csv"
	// var sqliteToCsvExpectedKeys = []string{
	// 	sqliteToCsvPrefix + "-source", sqliteToCsvPrefix + "-destination",
	// }

	// err := hasExpectedKeys(etlP, sqliteToCsvExpectedKeys)
	// if err != nil {
	// 	return err
	// }

	for _, table := range constants.ITSLOG_TABLES {
		query := "SELECT * FROM " + table

		// The CSV file will be based on the app name, the date, and
		// the table name. It will be stored to the same path as the SQLite
		// databases.
		csvFilename := strings.Join(
			[]string{
				etlP.AppId,
				etlP.Storage.YYYYMMDD(),
				table},
			"_") + ".csv"
		csvFullPath := path.Join(viper.GetString("storage.path"), csvFilename)

		err := queryToCSV(etlP.Storage, query, table, csvFullPath)
		if err != nil {
			return err
		}

	}
	return nil
}
