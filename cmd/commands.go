package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/xiayesuifeng/gopanel/core/config"
	"log"
	"os"
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
	err := config.ParseConf(cfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("please config config.json")
		}

		log.Fatalln(err)
	}
}
