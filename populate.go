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
	typeName := t.Name()

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

		// Create a unique key for field-level sequence management
		sequenceKey := typeName + "." + fieldType.Name
		value := generateFromTag(tag, field.Type(), sequenceKey)

		if value != nil {
			setFieldValue(field, value)
		}
	}
}

// tagContext holds context for tag generation.
type tagContext struct {
	parts       []string
	fieldType   reflect.Type
	sequenceKey string
}

// simpleTagGenerators maps tag names to simple generators (no parameters needed).
var simpleTagGenerators = map[string]func() any{
	"name":      func() any { return gofakeit.Name() },
	"email":     func() any { return gofakeit.Email() },
	"phone":     func() any { return gofakeit.Phone() },
	"company":   func() any { return gofakeit.Company() },
	"word":      func() any { return gofakeit.Word() },
	"uuid":      func() any { return gofakeit.UUID() },
	"url":       func() any { return gofakeit.URL() },
	"username":  func() any { return gofakeit.Username() },
	"datetime":  func() any { return gofakeit.Date() },
	"firstname": func() any { return gofakeit.FirstName() },
	"lastname":  func() any { return gofakeit.LastName() },
	"city":      func() any { return gofakeit.City() },
	"country":   func() any { return gofakeit.Country() },
	"zip":       func() any { return gofakeit.Zip() },
	"latitude":  func() any { return gofakeit.Latitude() },
	"longitude": func() any { return gofakeit.Longitude() },
	"ipv4":      func() any { return gofakeit.IPv4Address() },
	"ipv6":      func() any { return gofakeit.IPv6Address() },
	"mac":       func() any { return gofakeit.MacAddress() },
	"hexcolor":  func() any { return gofakeit.HexColor() },
	"rgbcolor":  func() any { return gofakeit.RGBColor() },
	"useragent": func() any { return gofakeit.UserAgent() },
	"address":   func() any { return gofakeit.Address().Address },
}

func generateFromTag(tag string, fieldType reflect.Type, sequenceKey string) any {
	parts := strings.Split(tag, ":")
	tagName := parts[0]

	// Check simple generators first
	if generator, ok := simpleTagGenerators[tagName]; ok {
		return generator()
	}

	// Handle complex generators that need context
	ctx := tagContext{
		parts:       parts,
		fieldType:   fieldType,
		sequenceKey: sequenceKey,
	}

	return handleComplexTag(tagName, ctx)
}

func handleComplexTag(tagName string, ctx tagContext) any {
	switch tagName {
	case "sentence":
		return handleSentenceTag(ctx.parts)
	case "range":
		return handleRangeTag(ctx.parts)
	case "sequence":
		return generateSequenceForField(ctx.parts, ctx.fieldType, ctx.sequenceKey)
	case "password":
		return handlePasswordTag(ctx.parts)
	case "date":
		return handleDateTag(ctx.parts)
	case "bool":
		return handleBoolTag(ctx.parts)
	case "oneof":
		return handleOneOfTag(ctx.parts)
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

// handlePasswordTag generates a password.
// Supports "password" for default (8-16 chars) or "password:N" for specific length.
func handlePasswordTag(parts []string) any {
	length := 12

	if len(parts) == 2 {
		if l, err := strconv.Atoi(parts[1]); err == nil && l > 0 {
			length = l
		}
	}

	return gofakeit.Password(true, true, true, true, false, length)
}

// handleDateTag generates a date string.
// Supports "date" for YYYY-MM-DD format or "date:format" for custom format.
func handleDateTag(parts []string) any {
	format := "2006-01-02"

	if len(parts) == 2 {
		format = parts[1]
	}

	return gofakeit.Date().Format(format)
}

// handleBoolTag generates a boolean value.
// Supports "bool" for random, "bool:true" for always true, "bool:false" for always false.
func handleBoolTag(parts []string) any {
	if len(parts) == 2 {
		switch strings.ToLower(parts[1]) {
		case "true":
			return true
		case "false":
			return false
		}
	}

	return gofakeit.Bool()
}

// handleOneOfTag selects a random value from the provided options.
// Usage: "oneof:a,b,c" returns one of "a", "b", or "c".
func handleOneOfTag(parts []string) any {
	if len(parts) != 2 {
		return ""
	}

	options := strings.Split(parts[1], ",")
	if len(options) == 0 {
		return ""
	}

	for i := range options {
		options[i] = strings.TrimSpace(options[i])
	}

	return options[gofakeit.Number(0, len(options)-1)]
}

// generateSequenceForField generates a sequence value using the global registry.
// Supports custom start values via "sequence:N" syntax (e.g., "sequence:100").
func generateSequenceForField(parts []string, fieldType reflect.Type, sequenceKey string) any {
	startValue := int64(1)

	if len(parts) == 2 {
		if start, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
			startValue = start
		}
	}

	counter := globalRegistry.getOrCreate(sequenceKey, startValue)
	next := counter.next()

	return convertSequenceValue(next, fieldType)
}

func convertSequenceValue(next int64, fieldType reflect.Type) any {
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
