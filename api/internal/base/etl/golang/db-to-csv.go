package etl

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

var sqliteToCsvPrefix string = "sqlite_to_csv"
var sqliteToCsvExpectedKeys = []string{
	sqliteToCsvPrefix + "_source", sqliteToCsvPrefix + "_destination",
}

func hasExpectedKeys(etlP *types.RunEtlParams, keylist []string) error {
	for _, key := range sqliteToCsvExpectedKeys {
		_, ok := etlP.Payload[key]
		if !ok {
			return fmt.Errorf("missing parameter %s for db-to-csv", key)
		}
	}
	return nil
}

func SqliteToCSV(etlP *types.RunEtlParams) error {
	// err := hasExpectedKeys(etlP, sqliteToCsvExpectedKeys)
	// if err != nil {
	// 	return err
	// }

	for _, table := range constants.ITSLOG_TABLES {
		// 1. Execute the query
		rows, err := etlP.Storage.GetDB().Query("SELECT * FROM " + table)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

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

		file, err := os.Create(csvFullPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		columns, err := rows.Columns()
		if err != nil {
			// FIXME: error handling
			log.Fatal(err)
		}
		if err := writer.Write(columns); err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePointers := make([]interface{}, len(columns))
			for i := range values {
				valuePointers[i] = &values[i]
			}

			if err := rows.Scan(valuePointers...); err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		log.Println("exported: " + csvFilename)
	}
	return nil
}
