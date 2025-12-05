package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfig() {
	// Environment variables for its-log want to be prefixed
	// with ITSLOG_ or we won't use them.
	viper.SetEnvPrefix("itslog")

	viper.AutomaticEnv()

	viper.SetConfigName("config")

	// The config file either needs to be in the
	// home directory, or next to the binary.
	viper.AddConfigPath("$HOME/.itslog")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	// If we can't find the config, the world is ending. Exit ungracefully.
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
