/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/jadudm/its-log/internal/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// summarizeCmd represents the summarize command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("summarize called")
		storage := &sqlite.SqliteStorage{
			Path: viper.GetString("sqlite.path"),
		}

		t := time.Now()
		// FIXME: Give the command line flexibility
		// to choose the date, or to go back N days, or...
		yesterday := t.AddDate(0, 0, -1)

		storage.Init(yesterday)
		storage.Summarize()
	},
}

func init() {
	rootCmd.AddCommand(summarizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// summarizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// summarizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
