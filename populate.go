package gofab

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

// autoPopulateFromTags automatically populates struct fields based on gofab tags
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
		
		// Skip unexported fields
		if !field.CanSet() {
			continue
		}
		
		// Get gofab tag
		tag := fieldType.Tag.Get("gofab")
		if tag == "" {
			// No tag - leave as zero value
			continue
		}
		
		if tag == "-" {
			// Skip this field
			continue
		}
		
		// Generate value based on tag
		value := generateFromTag(tag, field.Type())
		if value != nil {
			setFieldValue(field, value)
		}
	}
}

// generateFromTag generates value based on gofab struct tag
func generateFromTag(tag string, fieldType reflect.Type) any {
	// Handle parameterized tags (e.g., "range:1,100", "sentence:3")
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
		count := 5 // default
		if len(parts) == 2 {
			if c, err := strconv.Atoi(parts[1]); err == nil && c > 0 {
				count = c
			}
		}
		return gofakeit.Sentence(count)
	case "range":
		if len(parts) == 2 {
			minMax := strings.Split(parts[1], ",")
			if len(minMax) == 2 {
				min, err1 := strconv.Atoi(strings.TrimSpace(minMax[0]))
				max, err2 := strconv.Atoi(strings.TrimSpace(minMax[1]))
				if err1 == nil && err2 == nil && min <= max {
					return gofakeit.Number(min, max)
				}
			}
		}
		// Fallback for invalid range
		return gofakeit.Number(1, 100)
	case "sequence":
		// Use existing sequence functionality
		return generateSequence(fieldType)
	default:
		// Unknown tag - return nil to leave as zero value
		return nil
	}
}

// sequenceCounters holds global sequence counters for different types
var sequenceCounters = make(map[reflect.Type]*sequenceCounter)

// generateSequence generates sequential values using existing sequence logic
func generateSequence(fieldType reflect.Type) any {
	counter, exists := sequenceCounters[fieldType]
	if !exists {
		counter = &sequenceCounter{value: 0}
		sequenceCounters[fieldType] = counter
	}
	
	next := counter.next()
	
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(next)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(next)
	case reflect.String:
		return strconv.FormatInt(next, 10)
	default:
		return int(next)
	}
}

// setFieldValue safely sets a field value using reflection
func setFieldValue(field reflect.Value, value any) {
	if !field.CanSet() {
		return
	}
	
	valueReflect := reflect.ValueOf(value)
	if valueReflect.Type().ConvertibleTo(field.Type()) {
		field.Set(valueReflect.Convert(field.Type()))
	}
}