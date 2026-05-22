package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/input"
	"github.com/yourorg/logslice/internal/output"
)

func main() {
	from := flag.String("from", "", "start of time range (RFC3339)")
	to := flag.String("to", "", "end of time range (RFC3339)")
	fields := flag.String("fields", "", "comma-separated field=value filters")
	regexField := flag.String("regex-field", "", "field name for regex filter")
	regexPat := flag.String("regex", "", "regex pattern for --regex-field")
	invert := flag.Bool("invert", false, "invert the combined filter result")
	limit := flag.Int("limit", 0, "max number of output lines (0 = unlimited)")
	outFile := flag.String("out", "", "output file path (default: stdout)")
	flag.Parse()

	files := flag.Args()

	reader, err := input.NewLineReader(files)
	if err != nil {
		log.Fatalf("input error: %v", err)
	}
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("read error: %v", err)
	}

	var filters []filter.Filter

	if *from != "" || *to != "" {
		tr, err := filter.ParseTimeRange(*from, *to)
		if err != nil {
			log.Fatalf("time range: %v", err)
		}
		filters = append(filters, tr)
	}

	if *fields != "" {
		mf, err := filter.NewMultiFieldFilter(strings.Split(*fields, ","))
		if err != nil {
			log.Fatalf("field filter: %v", err)
		}
		filters = append(filters, mf)
	}

	if *regexField != "" || *regexPat != "" {
		rf, err := filter.NewRegexFilter(*regexField, *regexPat)
		if err != nil {
			log.Fatalf("regex filter: %v", err)
		}
		filters = append(filters, rf)
	}

	if *limit > 0 {
		lf, err := filter.NewLimitFilter(*limit)
		if err != nil {
			log.Fatalf("limit filter: %v", err)
		}
		filters = append(filters, lf)
	}

	var combined filter.Filter
	if len(filters) == 0 {
		combined = filter.PassAll{}
	} else {
		combined = filter.NewMultiFilter(filters)
	}
	if *invert {
		combined, err = filter.NewInvertFilter(combined)
		if err != nil {
			log.Fatalf("invert filter: %v", err)
		}
	}

	pipeline := filter.NewPipeline(lines, combined)
	results, err := pipeline.Run()
	if err != nil {
		log.Fatalf("pipeline: %v", err)
	}

	w, err := output.NewWriter(*outFile)
	if err != nil {
		log.Fatalf("output: %v", err)
	}
	defer w.Close()

	for _, line := range results {
		if err := w.WriteLine(line); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		}
	}
}
