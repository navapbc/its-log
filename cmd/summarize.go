/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/jadudm/its-log/internal/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var filename *string

// summarizeCmd represents the summarize command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarize the data in a database",
	Long:  `Provide the name of a database at sqlite.path to be summarized.`,
	Run: func(cmd *cobra.Command, args []string) {
		fullPath := filepath.Join(viper.GetString("sqlite.path"), *filename)
		fmt.Printf("Summarizing `%s`\n", fullPath)
		storage := &sqlite.SqliteStorage{
			Path: fullPath,
		}

		storage.Init()
		storage.Summarize()
	},
}

func init() {
	pfix := time.Now().Format("2006-01-02")
	filename = summarizeCmd.Flags().StringP("filename", "f", pfix+".sqlite", "name of database in sqlite.path")
	rootCmd.AddCommand(summarizeCmd)
}
