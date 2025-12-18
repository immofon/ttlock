package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var helloCmd = &cli.Command{
	Name:  "hello",
	Usage: "print hello",
	Action: func(c *cli.Context) error {
		fmt.Println("hello")
		return nil
	},
}
