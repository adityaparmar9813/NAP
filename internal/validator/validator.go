package validator

import (
	"fmt"
	"reflect"

	"github.com/adityaparmar9813/NAP/internal/types"
)

type ValidatorInterface interface {
	ValidateType(value interface{}, fieldType types.FieldType) error
}

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateType(value interface{}, fieldType types.FieldType) error {
	switch fieldType {
	case types.TypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %v", reflect.TypeOf(value))
		}
	case types.TypeInt:
		if _, ok := value.(int); !ok {
			return fmt.Errorf("expected int, got %v", reflect.TypeOf(value))
		}
	case types.TypeFloat:
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("expected float, got %v", reflect.TypeOf(value))
		}
	case types.TypeBoolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %v", reflect.TypeOf(value))
		}
	default:
		return fmt.Errorf("unknown field type: %s", fieldType)
	}
	return nil
}

func MatchesCriteria(record, criteria map[string]interface{}) bool {
	for key, criteriaValue := range criteria {
		recordValue, exists := record[key]
		if !exists {
			return false
		}

		if !compareValues(recordValue, criteriaValue) {
			return false
		}
	}
	return true
}

func compareValues(v1, v2 interface{}) bool {
	rv1 := reflect.ValueOf(v1)
	rv2 := reflect.ValueOf(v2)

	// Handle nil values
	if !rv1.IsValid() || !rv2.IsValid() {
		return rv1.IsValid() == rv2.IsValid()
	}

	// If types are the same, use direct comparison
	if rv1.Type() == rv2.Type() {
		return reflect.DeepEqual(v1, v2)
	}

	// Handle numeric comparisons
	switch rv1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch rv2.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return rv1.Int() == rv2.Int()
		case reflect.Float32, reflect.Float64:
			return float64(rv1.Int()) == rv2.Float()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch rv2.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return rv1.Uint() == rv2.Uint()
		case reflect.Float32, reflect.Float64:
			return float64(rv1.Uint()) == rv2.Float()
		}
	case reflect.Float32, reflect.Float64:
		switch rv2.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return rv1.Float() == float64(rv2.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return rv1.Float() == float64(rv2.Uint())
		case reflect.Float32, reflect.Float64:
			return rv1.Float() == rv2.Float()
		}
	}

	// For other types, try string comparison as a last resort
	return fmt.Sprintf("%v", v1) == fmt.Sprintf("%v", v2)
}
