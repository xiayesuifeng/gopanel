/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	port int
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Run:   serverRun,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "server listen port")
}

func serverRun(cmd *cobra.Command, args []string) {
	instance, err := core.New(port)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		if err := instance.Start(cmd.Context()); err != nil {
			log.Fatalln(err)
		}
	}()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		if err := instance.Close(); err != nil {
			log.Fatalln(err)
		}
	}
}
