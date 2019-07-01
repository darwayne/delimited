package delimited

import (
	"encoding/csv"
	"io"
)

type Writer struct {
	*csv.Writer
}

type WriterOptsFn func(reader *csv.Writer)

func TabWriterOpt() WriterOptsFn {
	return func(writer *csv.Writer) {
		writer.Comma = '\t'
	}
}

func CommaWriterOpt() WriterOptsFn {
	return func(writer *csv.Writer) {
		writer.Comma = ','
	}
}

func SpaceWriterOpt() WriterOptsFn {
	return func(writer *csv.Writer) {
		writer.Comma = ' '
	}
}

func NewWriter(writer io.Writer, opts ...WriterOptsFn) *Writer {
	csvWriter := csv.NewWriter(writer)
	result := &Writer{
		Writer: csvWriter,
	}

	for _, opt := range opts {
		opt(csvWriter)
	}

	return result
}
