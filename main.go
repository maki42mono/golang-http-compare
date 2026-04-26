package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"speedtest/parser"
	"syscall"
)

var registry = map[string]parser.TestCase{
	"timer": parser.TimerTest,
	"sync":  parser.SyncCase,
}

func main() {
	caseName := flag.String("case", "timer", "Test case name")
	timeout := flag.Duration("timeout", 0, "give just a timeout in ms")
	flag.Parse()

	runCase, ok := registry[*caseName]

	if !ok {
		fmt.Printf("unknown case: %v\n", *caseName)
		os.Exit(1)
	}

	fmt.Printf("i have such a case: %v\n", *caseName)

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	var links []string
	links = append(links, "https://serebii.net")

	result, error := runCase(ctx, links, 2)
	if error == context.Canceled {
		fmt.Printf("Oppsi poopsi. Why did you terminate?))\n")
		os.Exit(1)
	}
	if error != nil {
		fmt.Printf("case has an error: %v\n", error.Error())
		os.Exit(1)
	}

	fmt.Printf("case finished: %v\n", result)
}
