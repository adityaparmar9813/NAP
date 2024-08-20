package schema

import (
	"fmt"

	"github.com/adityaparmar9813/NAP/internal/types"
	"github.com/adityaparmar9813/NAP/internal/validator"
)

// Field represents a schema field
type Field struct {
	Name      string
	Type      types.FieldType
	Required  bool
	Validator types.ValidatorFunc
}

// Schema represents the document schema
type Schema struct {
	Fields map[string]Field
}

// NewSchema creates a new Schema
func NewSchema() *Schema {
	return &Schema{
		Fields: make(map[string]Field),
	}
}

// AddField adds a new field to the schema
func (s *Schema) AddField(name string, fieldType types.FieldType, required bool, validator types.ValidatorFunc) error {
	if _, exists := s.Fields[name]; exists {
		return fmt.Errorf("field '%s' already exists in schema", name)
	}

	s.Fields[name] = Field{
		Name:      name,
		Type:      fieldType,
		Required:  required,
		Validator: validator,
	}

	return nil
}

// Validate checks if a document conforms to the schema
func (s *Schema) Validate(doc map[string]interface{}) error {
	for fieldName, field := range s.Fields {
		value, exists := doc[fieldName]

		if !exists {
			if field.Required {
				return fmt.Errorf("required field '%s' is missing", fieldName)
			}
			continue // Skip validation for non-required fields that are not present
		}

		if err := validator.ValidateType(value, field.Type); err != nil {
			return fmt.Errorf("field '%s': %w", fieldName, err)
		}

		if field.Validator != nil && !field.Validator(value) {
			return fmt.Errorf("field '%s': failed custom validation", fieldName)
		}
	}

	return nil
}
