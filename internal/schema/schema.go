package schema

import (
	"fmt"

	"github.com/adityaparmar9813/NAP/internal/types"
)

// Field represents a schema field
type Field struct {
	Name     string
	Type     types.FieldType
	Required bool
}

// SchemaInterface defines the methods for a schema
type SchemaInterface interface {
	AddField(field Field) error
	Validate(doc map[string]interface{}) error
	PrintSchema()
	AddRecord(doc map[string]interface{}, storage StorageInterface) error
}

// StorageInterface defines the methods for storage operations
type StorageInterface interface {
	SaveStructToFile(data interface{}, filepath string) error
	LoadStructFromFile(filepath string, data interface{}) error
}

// ValidatorInterface defines the methods for validation
type ValidatorInterface interface {
	ValidateType(value interface{}, fieldType types.FieldType) error
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
func (s *Schema) AddField(field Field) error {
	if _, exists := s.Fields[field.Name]; exists {
		return fmt.Errorf("field '%s' already exists in schema", field.Name)
	}

	s.Fields[field.Name] = field

	return nil
}

// Validate checks if a document conforms to the schema
func (s *Schema) Validate(doc map[string]interface{}, validator ValidatorInterface) error {
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
	}

	return nil
}

// PrintSchema prints the schema
func (s *Schema) PrintSchema() {
	fmt.Println("Schema:")
	for name, field := range s.Fields {
		fmt.Printf("%s (%s, required=%t)\n", name, field.Type, field.Required)
	}
}

// AddRecord validates and saves a document to storage
func (s *Schema) AddRecord(doc map[string]interface{}, validator ValidatorInterface, storage StorageInterface) error {
	err := s.Validate(doc, validator)
	if err != nil {
		return err
	}

	err = storage.SaveStructToFile(doc, "./collections/user.json")
	if err != nil {
		return err
	}

	return nil
}

// BuildSchema creates a new schema and adds fields to it
func BuildSchema(storage StorageInterface, fields ...Field) (*Schema, error) {
	schema := NewSchema()

	for _, field := range fields {
		err := schema.AddField(field)
		if err != nil {
			return nil, err
		}
	}

	// Save the schema to a file
	err := storage.SaveStructToFile(schema, "./schemas/user.json")
	if err != nil {
		return nil, err
	}

	return schema, nil
}

// Test function to demonstrate usage
func Test(storage StorageInterface, validator ValidatorInterface) {
	nameField := Field{
		Name:     "name",
		Type:     types.TypeString,
		Required: true,
	}

	ageField := Field{
		Name:     "age",
		Type:     types.TypeInt,
		Required: true,
	}

	emailField := Field{
		Name:     "email",
		Type:     types.TypeString,
		Required: false,
	}

	schema, err := BuildSchema(storage, nameField, ageField, emailField)
	if err != nil {
		fmt.Println("Error adding fields:", err)
		return
	}

	schema.PrintSchema()

	doc := map[string]interface{}{
		"name":  "Arpit Dubey",
		"age":   22,
		"email": "adubey_be21@thapar.edu",
	}

	err = storage.LoadStructFromFile("./schemas/user.json", schema)
	if err != nil {
		fmt.Println("Error loading schema:", err)
		return
	}

	schema.PrintSchema()

	err = schema.AddRecord(doc, validator, storage)
	if err != nil {
		fmt.Println("Error adding record:", err)
		return
	}
	fmt.Println("Record added successfully")

	// Validate the document against the schema
	err = schema.Validate(doc, validator)
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
		err = schema.Validate(doc, validator)
		if err != nil {
			fmt.Printf("Invalid document %d: %v\n", i+1, err)
		}
	}
}
