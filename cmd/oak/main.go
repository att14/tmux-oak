package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	session := flag.String("session", "", "tmux session name")
	flag.Parse()

	if *showVersion {
		fmt.Printf("oak %s (%s)\n", version, commit)
		os.Exit(0)
	}

	if *session == "" {
		fmt.Fprintln(os.Stderr, "oak: --session is required")
		os.Exit(1)
	}

	fmt.Printf("oak %s — session: %s (TUI not yet implemented)\n", version, *session)
}
