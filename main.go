package main

import (
	"log"
	"os"

	"github.com/dll-as/gitc/cmd"
)

func main() {
	if err := cmd.Commands.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
