package excel

import (
	"fmt"
	"io"
	"reflect"
)

func Unmarshal(rdr Reader, target any) (*ImportResult, error) {
	policy, err := NewPolicy(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create Policy: %w", err)
	}

	headers := rdr.Headers()
	headerToField := make(map[string]string)
	for _, header := range headers {
		if fieldName, ok := policy.Mapping[header]; ok {
			headerToField[header] = fieldName
		}
	}

	var validationErrors []ValidationError
	rowNum := policy.HeadingRow
	newSlice := reflect.MakeSlice(policy.sliceVal.Type(), 0, 0)

	for {
		rowNum++
		row, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Row:     rowNum,
				Message: err.Error(),
			})
			continue
		}

		data := make(map[string]any, len(headerToField))
		for header, fieldName := range headerToField {
			if val, ok := row[header]; ok {
				data[fieldName] = val
			}
		}

		validation := Validate(policy, data)
		if validation.IsFail() {
			for field, errs := range validation.Errors {
				for rule, msg := range errs {
					validationErrors = append(validationErrors, ValidationError{
						Row:     rowNum,
						Field:   field,
						Rule:    rule,
						Message: msg,
					})
				}
			}
			continue
		}

		newStructPtr := reflect.New(policy.elemType)
		if err := validation.BindSafeData(newStructPtr.Interface()); err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Row:     rowNum,
				Message: fmt.Sprintf("failed to bind data: %v", err),
			})
			continue
		}

		if policy.isPointerElem {
			newSlice = reflect.Append(newSlice, newStructPtr)
		} else {
			newSlice = reflect.Append(newSlice, newStructPtr.Elem())
		}
	}

	policy.sliceVal.Set(newSlice)

	return &ImportResult{
		ValidationErrors: validationErrors,
	}, nil
}
