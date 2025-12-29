/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"os"

	"github.com/jadudm/its-log/internal/etl"
	"github.com/spf13/cobra"
)

type etlParamsT struct {
	etlRunscript *string
	etlFilename  *string
}

var etlParams etlParamsT

// etlCmd represents the etl command
var etlCmd = &cobra.Command{
	Use:   "etl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite", *etlParams.etlFilename)
		// FIXME: Should I create this if it doesn't exist?
		if err != nil {
			panic(err)
		}

		script, err := os.ReadFile(*etlParams.etlRunscript)
		if err != nil {
			panic(err)
		}
		etl.Run(string(script), db)
	},
}

func init() {
	rootCmd.AddCommand(etlCmd)

	// etlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// This should take a path to an ETL "script" (JSonnet? Json?)
	// and a path to an SQLite DB
	etlParams.etlRunscript = etlCmd.Flags().StringP("runscript", "r", "REQUIRED", "path to runscript")
	etlParams.etlFilename = etlCmd.Flags().StringP("sqlite", "s", "REQUIRED", "path to SQLite file")
	etlCmd.MarkFlagRequired("runscript")
	etlCmd.MarkFlagRequired("sqlite")
}
