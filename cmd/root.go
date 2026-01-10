package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "its-log",
		Short: "A simple, valueless event logger",
		Long: `It's better than bad, it's good!
	
its-log is a server and command-line tool for collecting and managing
simple, event-based data about applications in resource-constrained,
NIST-controlled environments. If you just want to count things, then
everyone wants its-log.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	// Everything we want in the env for Viper
	// must be prefixed with ITSLOG_
	viper.SetEnvPrefix("ITSLOG")
	viper.AutomaticEnv()
}
