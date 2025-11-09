package excel

import (
	"io"
)

type ImportResult struct {
	ValidationErrors []ValidationError
}

type ValidationError struct {
	Row     int
	Field   string
	Rule    string
	Message string
}

func Import(r io.Reader, target any) (*ImportResult, error) {
	reader, err := NewCsvReader(r)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return Unmarshal(reader, target)
}
