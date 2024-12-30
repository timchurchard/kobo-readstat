package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/timchurchard/kobo-readstat/cmd"
)

const cliName = "kobo-readstat"

func main() {
	//if len(os.Args) < 2 {
	//	os.Exit(cmd.Gui(os.Stdout))
	//}

	// Save the command and reset the flags
	command := os.Args[1]
	flag.CommandLine = flag.NewFlagSet(cliName, flag.ExitOnError)
	os.Args = append([]string{cliName}, os.Args[2:]...)

	switch command {
	case "sync":
		os.Exit(cmd.Sync(os.Stdout))

	case "stats":
		os.Exit(cmd.Stats(os.Stdout))

	//case "gui":
	//	os.Exit(cmd.Gui(os.Stdout))

	default:
		usageRoot()
	}
}

func usageRoot() {
	fmt.Printf("usage: %s commands(gui, sync or stats) options\n", cliName)
	os.Exit(1)
}
