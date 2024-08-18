package main

import (
	"fmt"
	"reflect"
)

// FieldType represents the type of a field
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeFloat   FieldType = "float"
	TypeBoolean FieldType = "boolean"
)

// ValidatorFunc is a function type for custom validators
type ValidatorFunc func(interface{}) bool

// Field represents a schema field
type Field struct {
	Name      string
	Type      FieldType
	Required  bool
	Validator ValidatorFunc
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
func (s *Schema) AddField(name string, fieldType FieldType, required bool, validator ValidatorFunc) error {
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

		if err := validateType(value, field.Type); err != nil {
			return fmt.Errorf("field '%s': %w", fieldName, err)
		}

		if field.Validator != nil && !field.Validator(value) {
			return fmt.Errorf("field '%s': failed custom validation", fieldName)
		}
	}

	return nil
}

func validateType(value interface{}, fieldType FieldType) error {
	switch fieldType {
	case TypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %v", reflect.TypeOf(value))
		}
	case TypeInt:
		if _, ok := value.(int); !ok {
			return fmt.Errorf("expected int, got %v", reflect.TypeOf(value))
		}
	case TypeFloat:
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("expected float, got %v", reflect.TypeOf(value))
		}
	case TypeBoolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %v", reflect.TypeOf(value))
		}
	default:
		return fmt.Errorf("unknown field type: %s", fieldType)
	}
	return nil
}

func main() {
	// Create a new schema
	schema := NewSchema()

	// Add fields to the schema
	err := schema.AddField("name", TypeString, true, func(v interface{}) bool {
		s, ok := v.(string)
		return ok && len(s) > 0
	})
	if err != nil {
		fmt.Println("Error adding field:", err)
		return
	}

	err = schema.AddField("age", TypeInt, true, func(v interface{}) bool {
		age, ok := v.(int)
		return ok && age >= 0 && age <= 120
	})
	if err != nil {
		fmt.Println("Error adding field:", err)
		return
	}

	err = schema.AddField("email", TypeString, false, func(v interface{}) bool {
		s, ok := v.(string)
		// This is a very simplistic email validation
		return ok && len(s) > 3 && s[0] != '@' && s[len(s)-1] != '@'
	})
	if err != nil {
		fmt.Println("Error adding field:", err)
		return
	}

	// Example valid document
	doc := map[string]interface{}{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}

	// Validate the document against the schema
	err = schema.Validate(doc)
	if err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("Document is valid")
	}

	// Invalid document examples
	invalidDocs := []map[string]interface{}{
		{"name": "", "age": 30},              // Invalid name (empty string)
		{"name": "Jane Doe", "age": 150},     // Invalid age (out of range)
		{"name": "Bob", "email": "@invalid"}, // Invalid email
		{"age": 25},                          // Missing required field (name)
	}

	for i, doc := range invalidDocs {
		err = schema.Validate(doc)
		if err != nil {
			fmt.Printf("Invalid document %d: %v\n", i+1, err)
		}
	}
}
