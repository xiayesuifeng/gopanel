package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"log"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:              "gopanel",
	PersistentPreRun: initConfig,
}

func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.json", "config file (default is config.json)")

	return rootCmd.Execute()
}

func initConfig(cmd *cobra.Command, args []string) {
	firstLaunch, err := config.ParseConf(cfgFile)

	if err != nil {
		log.Fatalln(err)
	}

	cmd.SetContext(context.WithValue(cmd.Context(), "firstLaunch", firstLaunch))
}
