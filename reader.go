package delimited

import (
	"context"
	"encoding/csv"
	"io"
)

type Reader struct {
	csvReader *csv.Reader
}

type ReaderOptsFn func(reader *csv.Reader)

func SpaceReaderOpt() ReaderOptsFn {
	return func(reader *csv.Reader) {
		reader.Comma = ' '
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
	}
}

func TabReaderOpt() ReaderOptsFn {
	return func(reader *csv.Reader) {
		reader.Comma = '\t'
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
	}
}

func CommaReaderOpt() ReaderOptsFn {
	return func(reader *csv.Reader) {
		reader.Comma = ','
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
	}
}

func NewReader(reader io.Reader, opts ...ReaderOptsFn) *Reader {
	csvReader := csv.NewReader(reader)
	result := &Reader{
		csvReader: csvReader,
	}

	for _, opt := range opts {
		opt(csvReader)
	}

	return result
}

func (r *Reader) EachRow(ctx context.Context, fn func([]string)) error {
	for {
		rec, err := r.csvReader.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fn(rec)
	}
}
