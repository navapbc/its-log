package cmd

import (
	"os"

	"github.com/jadudm/its-log/internal/etl"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

var etlParams etl.EtlParams

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

		script, err := os.ReadFile(etlParams.EtlRunscriptPath)
		if err != nil {
			panic(err)
		}
		etlParams.EtlApiKey = os.Getenv("ITSLOG_APIKEY")
		etlParams.EtlUrl = gjson.Get(string(script), "server.url").String()
		etl.Run(string(script), etlParams)
	},
}

func init() {
	rootCmd.AddCommand(etlCmd)

	// etlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// This should take a path to an ETL "script" (JSonnet? Json?)
	// and a path to an SQLite DB
	etlCmd.Flags().StringVar(&etlParams.EtlRunscriptPath, "runscript", "REQUIRED", "path to runscript")
	etlCmd.Flags().StringVar(&etlParams.EtlDate, "date", "REQUIRED", "date to run ETL on")
	etlCmd.MarkFlagRequired("runscript")
	etlCmd.MarkFlagRequired("date")
}
