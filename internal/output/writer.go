package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Writer writes filtered log lines to a destination.
type Writer struct {
	w   *bufio.Writer
	owned bool // true if we opened the file and must close it
	f   *os.File
}

// NewWriter creates a Writer targeting the given path.
// If path is empty, stdout is used.
func NewWriter(path string) (*Writer, error) {
	if path == "" {
		return &Writer{w: bufio.NewWriter(os.Stdout)}, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("output: cannot open %q: %w", path, err)
	}
	return &Writer{w: bufio.NewWriter(f), f: f, owned: true}, nil
}

// NewWriterFromIO wraps an existing io.Writer (useful for tests).
func NewWriterFromIO(dst io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(dst)}
}

// WriteLine writes a single line followed by a newline character.
func (w *Writer) WriteLine(line string) error {
	_, err := fmt.Fprintln(w.w, line)
	return err
}

// Flush flushes buffered output to the underlying writer.
func (w *Writer) Flush() error {
	return w.w.Flush()
}

// Close flushes and closes the underlying file if owned.
func (w *Writer) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	if w.owned && w.f != nil {
		return w.f.Close()
	}
	return nil
}
