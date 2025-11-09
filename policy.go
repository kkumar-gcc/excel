package excel

import (
	"reflect"
	"strings"
)

const (
	tagExcel    = "excel"
	tagCSV      = "csv"
	tagValidate = "validate"
	tagLabel    = "label"
	tagMessage  = "message"
	tagSkip     = "-"
)

type Ruler interface {
	Rules() map[string]string
}

type Filterer interface {
	Filters() map[string]string
}

type Attributer interface {
	Attributes() map[string]string
}

type Messager interface {
	Messages() map[string]string
}

type HeadingRowProvider interface {
	HeadingRow() int
}

type Policy struct {
	Mapping    map[string]string
	Rules      map[string]string
	Filters    map[string]string
	Attributes map[string]string
	Messages   map[string]string
	HeadingRow int

	elemType      reflect.Type
	sliceVal      reflect.Value
	isPointerElem bool
	fieldCache    map[string]reflect.StructField
}

func NewPolicy(target any) (*Policy, error) {
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.Elem().Kind() != reflect.Slice {
		return nil, ErrInvalidTarget
	}

	sliceVal := targetVal.Elem()
	sliceType := sliceVal.Type()
	sliceElemType := sliceType.Elem()
	isPointerElem := sliceElemType.Kind() == reflect.Ptr

	elemType := sliceElemType
	if isPointerElem {
		elemType = sliceElemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		return nil, ErrTargetMustBeSliceOfStructs
	}

	policy := &Policy{
		elemType:      elemType,
		sliceVal:      sliceVal,
		isPointerElem: isPointerElem,
		Mapping:       make(map[string]string),
		Rules:         make(map[string]string),
		Filters:       make(map[string]string),
		Attributes:    make(map[string]string),
		Messages:      make(map[string]string),
		HeadingRow:    1,
		fieldCache:    make(map[string]reflect.StructField, elemType.NumField()),
	}

	policy.parseStructTags()

	instance := reflect.New(elemType).Interface()
	if v, ok := instance.(Ruler); ok {
		for k, v := range v.Rules() {
			policy.Rules[k] = v
		}
	}
	if v, ok := instance.(Filterer); ok {
		for k, v := range v.Filters() {
			policy.Filters[k] = v
		}
	}
	if v, ok := instance.(Attributer); ok {
		for k, v := range v.Attributes() {
			policy.Attributes[k] = v
		}
	}
	if v, ok := instance.(Messager); ok {
		for k, v := range v.Messages() {
			policy.Messages[k] = v
		}
	}
	if v, ok := instance.(HeadingRowProvider); ok {
		policy.HeadingRow = v.HeadingRow()
	}

	return policy, nil
}

func (r *Policy) parseStructTags() {
	for i := 0; i < r.elemType.NumField(); i++ {
		field := r.elemType.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		r.fieldCache[fieldName] = field

		tag := field.Tag.Get(tagExcel)
		if tag == "" {
			tag = field.Tag.Get(tagCSV)
		}
		if tag == tagSkip {
			continue
		}

		headerName := tag
		if headerName == "" {
			headerName = fieldName
		}

		for _, header := range strings.Split(headerName, ",") {
			if header = strings.TrimSpace(header); header != "" {
				r.Mapping[header] = fieldName
			}
		}

		if tag, ok := field.Tag.Lookup(tagValidate); ok {
			r.Rules[fieldName] = tag
		}
		if tag, ok := field.Tag.Lookup(tagLabel); ok {
			r.Attributes[fieldName] = tag
		}
		if tag, ok := field.Tag.Lookup(tagMessage); ok {
			for _, pair := range strings.Split(tag, "|") {
				if parts := strings.SplitN(pair, ":", 2); len(parts) == 2 {
					if rule, msg := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]); rule != "" && msg != "" {
						r.Messages[fieldName+"."+rule] = msg
					}
				}
			}
		}
	}
}

func (r *Policy) FieldByName(name string) (reflect.StructField, bool) {
	field, ok := r.fieldCache[name]
	return field, ok
}
