package excel

import (
	"encoding/csv"
	"io"
)

type CsvReaderOption func(*CsvReader)

func WithCsvHeadingRow(row int) CsvReaderOption {
	return func(r *CsvReader) {
		if row > 0 {
			r.headingRow = row
		}
	}
}

func WithCsvComma(comma rune) CsvReaderOption {
	return func(r *CsvReader) {
		r.reader.Comma = comma
	}
}

type CsvReader struct {
	reader  *csv.Reader
	headers []string
	closer  io.Closer

	headingRow int
}

func NewCsvReader(rdr io.Reader, opts ...CsvReaderOption) (*CsvReader, error) {
	var closer io.Closer
	if c, ok := rdr.(io.Closer); ok {
		closer = c
	}

	csvReader := csv.NewReader(rdr)
	csvReader.TrimLeadingSpace = true

	reader := &CsvReader{
		reader:     csvReader,
		headingRow: 1,
		closer:     closer,
	}

	for _, opt := range opts {
		opt(reader)
	}

	for i := 1; i < reader.headingRow; i++ {
		_, err := reader.reader.Read()
		if err != nil {
			if err == io.EOF {
				return nil, ErrNoHeaders
			}
			return nil, err
		}
	}

	headers, err := reader.reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, ErrNoHeaders
		}
		return nil, err
	}

	if len(headers) == 0 {
		return nil, ErrNoHeaders
	}

	reader.headers = headers
	return reader, nil
}

func (r *CsvReader) Close() error {
	if r.closer != nil {
		return r.closer.Close()
	}
	return nil
}

func (r *CsvReader) Headers() []string {
	return r.headers
}

func (r *CsvReader) Read() (map[string]string, error) {
	record, err := r.reader.Read()
	if err != nil {
		return nil, err
	}

	row := make(map[string]string, len(r.headers))
	for i, header := range r.headers {
		if i < len(record) {
			row[header] = record[i]
		} else {
			row[header] = ""
		}
	}
	return row, nil
}
