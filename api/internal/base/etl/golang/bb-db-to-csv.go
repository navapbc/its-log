package etl

import (
	"path"
	"strings"

	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

var BB_CSV_OUTPUTS = map[string]map[string]string{
	"itslog_summary": {
		"global":    "tags IS NULL",
		"app_level": "tags IS NOT NULL",
	},
	"itslog_etl":    {"": "*"},
	"itslog_events": {"": "*"},
}

func BBSqliteToCSV(etlP *types.RunEtlParams) error {
	for table, outputs := range BB_CSV_OUTPUTS {
		for file_suffix, filter := range outputs {
			// 1. Execute the query
			query := "SELECT * FROM " + table
			if filter != "*" {
				query += " WHERE " + filter
			}

			// The CSV file will be based on the app name, the date, and
			// the table name. It will be stored to the same path as the SQLite
			// databases.
			csvFilename := strings.Join(
				[]string{
					etlP.AppId,
					etlP.Storage.YYYYMMDD(),
					table},
				"_")
			if file_suffix != "" {
				csvFilename += "_" + file_suffix
			}
			csvFilename += ".csv"
			csvFullPath := path.Join(viper.GetString("storage.path"), csvFilename)

			err := queryToCSV(etlP.Storage, query, table, csvFullPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
