package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/ansod/nhl-terminal/helpers"
)

const (
	name = "nhl-terminal"
	desc = "" // TODO: Add description
)

func main() {

	parser := argparse.NewParser(name, desc)

	var args []bool

	goals := parser.Flag("g", "goal", &argparse.Options{Help: "Get alerted when a goal has been scored."})
	start := parser.Flag("s", "start", &argparse.Options{Help: "Get alerted when a game starts."})
	end := parser.Flag("e", "end", &argparse.Options{Help: "Get alerted when a game has ended."})

	args = append(args, *goals, *start, *end)

	err := parser.Parse(os.Args)

	if err != nil {
		// handle error
		fmt.Println(err)
	}

	update(args)

	response := helpers.Get()

}

func update(flags []bool) {

}
