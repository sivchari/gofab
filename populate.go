package gofab

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

func autoPopulateFromTags(obj any) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		tag := fieldType.Tag.Get("gofab")
		if tag == "" {
			continue
		}

		if tag == "-" {
			continue
		}

		value := generateFromTag(tag, field.Type())
		if value != nil {
			setFieldValue(field, value)
		}
	}
}

func generateFromTag(tag string, fieldType reflect.Type) any {
	return handleTagGeneration(tag, fieldType)
}

func handleTagGeneration(tag string, fieldType reflect.Type) any {
	parts := strings.Split(tag, ":")
	tagName := parts[0]

	switch tagName {
	case "name":
		return gofakeit.Name()
	case "email":
		return gofakeit.Email()
	case "phone":
		return gofakeit.Phone()
	case "company":
		return gofakeit.Company()
	case "address":
		return gofakeit.Address().Address
	case "word":
		return gofakeit.Word()
	case "sentence":
		return handleSentenceTag(parts)
	case "range":
		return handleRangeTag(parts)
	case "sequence":
		return generateSequence(fieldType)
	default:
		return nil
	}
}

func handleSentenceTag(parts []string) any {
	count := 5

	if len(parts) == 2 {
		if c, err := strconv.Atoi(parts[1]); err == nil && c > 0 {
			count = c
		}
	}

	return gofakeit.Sentence(count)
}

func handleRangeTag(parts []string) any {
	if len(parts) != 2 {
		return gofakeit.Number(1, 100)
	}

	minMax := strings.Split(parts[1], ",")
	if len(minMax) != 2 {
		return gofakeit.Number(1, 100)
	}

	minVal, err1 := strconv.Atoi(strings.TrimSpace(minMax[0]))
	maxVal, err2 := strconv.Atoi(strings.TrimSpace(minMax[1]))

	if err1 == nil && err2 == nil && minVal <= maxVal {
		return gofakeit.Number(minVal, maxVal)
	}

	return gofakeit.Number(1, 100)
}

var sequenceCounters = make(map[reflect.Type]*sequenceCounter)

func generateSequence(fieldType reflect.Type) any {
	counter, exists := sequenceCounters[fieldType]
	if !exists {
		counter = &sequenceCounter{value: 0}
		sequenceCounters[fieldType] = counter
	}

	next := counter.next()

	switch fieldType.Kind() { //nolint:exhaustive // exhaustive switch is not necessary here
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(next)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if next < 0 {
			return uint(0)
		}

		return uint(next)
	case reflect.String:
		return strconv.FormatInt(next, 10)
	case reflect.Bool:
		return next%2 == 0
	case reflect.Float32:
		return float32(next)
	case reflect.Float64:
		return float64(next)
	default:
		return int(next)
	}
}

func setFieldValue(field reflect.Value, value any) {
	if !field.CanSet() {
		return
	}

	valueReflect := reflect.ValueOf(value)
	if valueReflect.Type().ConvertibleTo(field.Type()) {
		field.Set(valueReflect.Convert(field.Type()))
	}
}
