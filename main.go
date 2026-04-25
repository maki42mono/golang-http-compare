package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"speedtest/parser"
)

var registry = map[string]parser.TestCase{
	"sync": parser.SyncCall,
}

func main() {
	caseName := flag.String("case", "sync", "Test case name")
	timeout := flag.Duration("timeout", 0, "give just a timeout in ms")
	flag.Parse()

	runCase, ok := registry[*caseName]

	if !ok {
		fmt.Printf("unknown case: %v\n", *caseName)
		os.Exit(1)
	}

	fmt.Printf("i have such a case: %v\n", *caseName)

	ctx := context.Background()

	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	result, error := runCase(ctx)
	if error != nil {
		fmt.Printf("case has an error: %v\n", error.Error())
		os.Exit(1)
	}

	fmt.Printf("case finished: %v\n", result)
}
