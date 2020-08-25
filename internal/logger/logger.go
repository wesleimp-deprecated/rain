package logger

import "github.com/apex/log"

// Writer writes with log.Info.
type Writer struct {
	entry *log.Entry
}

// NewWriter creates a new log writer.
func NewWriter(entry *log.Entry) Writer {
	return Writer{entry}
}

func (t Writer) Write(p []byte) (n int, err error) {
	t.entry.Info(string(p))
	return len(p), nil
}

// ErrorWriter writes with log.Error.
type ErrorWriter struct {
	entry *log.Entry
}

// NewErrWriter creates a new log writer.
func NewErrWriter(entry *log.Entry) ErrorWriter {
	return ErrorWriter{entry}
}

func (w ErrorWriter) Write(p []byte) (n int, err error) {
	w.entry.Error(string(p))
	return len(p), nil
}
