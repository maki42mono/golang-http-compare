package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"speedtest/parser"
	"syscall"
	"time"
)

var registry = map[string]parser.TestCase{
	// "timer": parser.TimerTest,
	"sync": parser.SyncCase,
}

func measure(w *csv.Writer) func() {
	start := time.Now()

	return func() {
		w.Write([]string{fmt.Sprintf("Time needed: %v", time.Since(start))})
	}
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
	file, err := os.Create(fmt.Sprintf("dump/res_%v.csv", *caseName))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	if *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}

	var links []string
	links = append(links, "https://lnb.lt")
	writer := csv.NewWriter(file)
	defer writer.Flush()
	defer measure(writer)()
	result, error := runCase(ctx, links, 2)
	if error == context.Canceled {
		fmt.Printf("Oppsi poopsi. Why did you terminate?))\n")
		os.Exit(1)
	}
	if error != nil {
		fmt.Printf("case has an error: %v\n", error.Error())
		os.Exit(1)
	}

	sorted := make([]*parser.ParsedLink, 0, len(result))
	for _, v := range result {
		sorted = append(sorted, v)
	}
	sort.Slice(sorted, func(i int, j int) bool {
		if sorted[i].Deep != sorted[j].Deep {
			return sorted[i].Deep > sorted[j].Deep
		}
		if sorted[i].Count != sorted[j].Count {
			return sorted[i].Count > sorted[j].Count
		}
		if sorted[i].OK != sorted[j].OK {
			return sorted[i].OK
		}
		return sorted[i].Link > sorted[j].Link
	})

	writer.Write([]string{"Link", "Depth", "Count", "Success"})
	fmt.Println("Started writing...")
	for _, v := range sorted {
		writer.Write(v.GetRow())
	}
	writer.Write([]string{fmt.Sprintf("Links parsed: %v", len(sorted))})
	fmt.Println("Done!")
}
