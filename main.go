package main

import (
	"fmt"
	"os"

	"github.com/jef/gh-vanity/cmd/ghvanity/commands"
)

func main() {
	command := commands.NewCommand()
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
