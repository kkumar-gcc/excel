package excel

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type ExcelReader struct {
	file    *excelize.File
	rows    *excelize.Rows
	headers []string

	sheetName  string
	headingRow int
}

type ExcelReaderOption func(*ExcelReader)

func WithSheetName(name string) ExcelReaderOption {
	return func(r *ExcelReader) {
		r.sheetName = name
	}
}

func WithHeadingRow(row int) ExcelReaderOption {
	return func(r *ExcelReader) {
		if row > 0 {
			r.headingRow = row
		}
	}
}

func NewExcelReader(rdr io.Reader, opts ...ExcelReaderOption) (excelReader *ExcelReader, err error) {
	f, err := excelize.OpenReader(rdr)
	if err != nil {
		return nil, err
	}

	excelReader = &ExcelReader{
		file:       f,
		headingRow: 1,
		sheetName:  f.GetSheetName(0),
	}

	defer func() {
		if err != nil {
			if excelReader.rows != nil {
				_ = excelReader.rows.Close()
			}
			if excelReader.file != nil {
				_ = excelReader.file.Close()
			}
		}
	}()

	for _, opt := range opts {
		opt(excelReader)
	}

	sheetName := excelReader.sheetName
	rows, err := f.Rows(sheetName)
	if err != nil {
		return nil, err
	}
	excelReader.rows = rows

	for i := 1; i < excelReader.headingRow; i++ {
		if !rows.Next() {
			return nil, io.EOF
		}
		_, _ = rows.Columns()
	}

	if !rows.Next() {
		return nil, io.EOF
	}

	headers, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return nil, ErrNoHeaders
	}

	excelReader.headers = headers

	return excelReader, nil
}

func (r *ExcelReader) Close() error {
	var firstErr error
	if r.rows != nil {
		if err := r.rows.Close(); err != nil {
			firstErr = err
		}
	}
	if r.file != nil {
		if err := r.file.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (r *ExcelReader) Headers() []string {
	return r.headers
}

func (r *ExcelReader) Read() (map[string]string, error) {
	if !r.rows.Next() {
		if err := r.rows.Error(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}

	cols, err := r.rows.Columns()
	if err != nil {
		return nil, err
	}

	record := make(map[string]string, len(r.headers))
	for i, h := range r.headers {
		if i < len(cols) {
			record[h] = cols[i]
		} else {
			record[h] = ""
		}
	}
	return record, nil
}
