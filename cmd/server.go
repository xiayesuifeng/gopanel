/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/xiayesuifeng/gopanel/api/server"
	"gitlab.com/xiayesuifeng/gopanel/core"
	"gitlab.com/xiayesuifeng/gopanel/web"
	"log"
	"strconv"
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

	if err := instance.Start(cmd.Context()); err != nil {
		log.Fatalln(err)
	}

	srv := server.NewServer(web.Assets())

	if err := srv.Run(":" + strconv.FormatInt(int64(port), 10)); err != nil {
		log.Fatalln(err)
	}
}
