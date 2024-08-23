package validator

import (
	"fmt"
	"reflect"

	"github.com/adityaparmar9813/NAP/internal/types"
)

// ValidatorInterface defines the methods for validating field types
type ValidatorInterface interface {
	ValidateType(value interface{}, fieldType types.FieldType) error
}

// Validator struct implements the ValidatorInterface
type Validator struct{}

// NewValidator creates a new Validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateType checks if the value conforms to the specified FieldType
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
