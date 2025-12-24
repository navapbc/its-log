/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path/filepath"
	"time"

	"github.com/jadudm/its-log/internal/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var filename string

// summarizeCmd represents the summarize command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarize the data in a database",
	Long:  `Provide the name of a database at sqlite.path to be summarized.`,
	Run: func(cmd *cobra.Command, args []string) {
		//
		storage := &sqlite.SqliteStorage{
			Path: filepath.Join(viper.GetString("sqlite.path"), filename),
		}

		storage.Init()
		storage.Summarize()
	},
}

func init() {
	pfix := time.Now().Format("2006-01-02")
	rootCmd.AddCommand(summarizeCmd)
	rootCmd.Flags().StringP("filename", "f", pfix+".sqlite", "name of database in sqlite.path")
}
