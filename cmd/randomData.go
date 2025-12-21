/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jadudm/its-log/internal/sqlite"
	"github.com/jadudm/its-log/internal/sqlite/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var row_count *int64
var app_count *int64
var events_per_app *int64

func randomTimeThisMonth() time.Time {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)
	duration := startOfNextMonth.Sub(startOfMonth)
	randomDuration := time.Duration(rand.Int63n(int64(duration)))
	randomTime := startOfMonth.Add(randomDuration)
	randomTime = randomTime.Truncate(time.Second)
	return randomTime
}

// randomDataCmd represents the randomData command
var randomDataCmd = &cobra.Command{
	Use:   "randomData",
	Short: "Generate authentic random data",
	Long:  `Generate random data.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("randomData called")
		storage := &sqlite.SqliteStorage{
			Path: viper.GetString("sqlite.path"),
		}
		storage.Init()

		for range *row_count {
			q := storage.GetQueries()
			// Do everything this month, randomly.
			//fmt.Sprintf("app_%04d", rand.Intn(int(*app_count+1)))
			the_time := randomTimeThisMonth()
			the_app := rand.Int63n(*app_count + 1)
			the_evt := rand.Int63n(*events_per_app + 1)
			// Insert a fake-stamped event from this month
			q.LogTimestampedEvent(context.Background(),
				models.LogTimestampedEventParams{
					Timestamp: the_time,
					Source:    the_app,
					Event:     the_evt,
				})
			// Insert a dictionary entry for this fake app/event
			q.UpdateDictionary(context.Background(),
				models.UpdateDictionaryParams{
					EventSource: fmt.Sprintf("app_%03d", the_app),
					EventName:   fmt.Sprintf("event_%03d", the_evt),
					SourceHash:  the_app,
					EventHash:   the_evt,
				})
		}

		storage.Close()
	},
}

func init() {
	rootCmd.AddCommand(randomDataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// randomDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	row_count = randomDataCmd.Flags().Int64("rows", 1000000, "number of random rows of data to generate")
	app_count = randomDataCmd.Flags().Int64("apps", 20, "number of simulated applications")
	events_per_app = randomDataCmd.Flags().Int64("events", 100, "number of possible events per simulated application")
}
