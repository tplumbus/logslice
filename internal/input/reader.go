package input

import (
	"bufio"
	"io"
	"os"
)

// LineReader reads lines from one or more sources sequentially.
type LineReader struct {
	sources []io.Reader
}

// NewLineReader creates a LineReader from the given file paths.
// If no paths are provided, stdin is used.
func NewLineReader(paths []string) (*LineReader, error) {
	if len(paths) == 0 {
		return &LineReader{sources: []io.Reader{os.Stdin}}, nil
	}

	sources := make([]io.Reader, 0, len(paths))
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		sources = append(sources, f)
	}
	return &LineReader{sources: sources}, nil
}

// Lines returns a channel that emits each line from all sources in order.
// The channel is closed when all sources are exhausted or ctx is done.
func (lr *LineReader) Lines() <-chan string {
	ch := make(chan string, 64)
	go func() {
		defer close(ch)
		for _, src := range lr.sources {
			scanner := bufio.NewScanner(src)
			for scanner.Scan() {
				line := scanner.Text()
				if line != "" {
					ch <- line
				}
			}
		}
	}()
	return ch
}
