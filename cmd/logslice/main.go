package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/input"
	"github.com/yourorg/logslice/internal/output"
)

func main() {
	var (
		from     = flag.String("from", "", "start of time range (RFC3339), inclusive")
		to       = flag.String("to", "", "end of time range (RFC3339), inclusive")
		field    = flag.String("field", "", "field filter in key=value format")
		outPath  = flag.String("out", "", "output file path (default: stdout)")
		timestampKey = flag.String("ts-key", "ts", "JSON key used for the timestamp field")
	)
	flag.Parse()

	// Build pipeline
	var opts []filter.Option

	if *from != "" || *to != "" {
		tr, err := filter.ParseTimeRange(*from, *to, *timestampKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid time range: %v\n", err)
			os.Exit(1)
		}
		opts = append(opts, filter.WithTimeRange(tr))
	}

	if *field != "" {
		fq, err := filter.ParseFieldQuery(*field)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid field query: %v\n", err)
			os.Exit(1)
		}
		opts = append(opts, filter.WithFieldQuery(fq))
	}

	pipeline := filter.NewPipeline(opts...)

	// Set up input
	lr, err := input.NewLineReader(flag.Args())
	if err != nil {
		log.Fatalf("error opening input: %v", err)
	}

	// Set up output
	w, err := output.NewWriter(*outPath)
	if err != nil {
		log.Fatalf("error opening output: %v", err)
	}
	defer w.Close()

	// Run
	for line := range pipeline.Run(lr.Lines()) {
		if err := w.WriteLine(line); err != nil {
			log.Fatalf("error writing output: %v", err)
		}
	}
}
