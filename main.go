package main

import (
	"gitlab.com/xiayesuifeng/gopanel/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
