package schema

import (
	"fmt"

	"github.com/adityaparmar9813/NAP/internal/storage"
	"github.com/adityaparmar9813/NAP/internal/types"
	"github.com/adityaparmar9813/NAP/internal/validator"
)

// Field represents a schema field
type Field struct {
	Name     string
	Type     types.FieldType
	Required bool
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
	}

	return nil
}

// Print a schema
func (s *Schema) PrintSchema() {
	fmt.Println("Schema:")
	for name, field := range s.Fields {
		fmt.Printf("%s (%s, required=%t)\n", name, field.Type, field.Required)
	}
}

func (s *Schema) AddRecord(doc map[string]interface{}) error {
	err := s.Validate(doc)
	if err != nil {
		return err
	}

	storage.SaveStructToFile(doc, "./collections/user.json")
	return nil
}

func Test() {
	// Create a new schema
	schema := NewSchema()

	// Create fields
	nameField := &Field{
		Name:     "name",
		Type:     types.TypeString,
		Required: true,
	}

	ageField := &Field{
		Name:     "age",
		Type:     types.TypeInt,
		Required: true,
	}

	emailField := &Field{
		Name:     "email",
		Type:     types.TypeString,
		Required: false,
	}

	// Add fields to the schema
	err := schema.AddField(*nameField)
	if err != nil {
		fmt.Println("Error adding name field:", err)
		return
	}

	err = schema.AddField(*ageField)
	if err != nil {
		fmt.Println("Error adding age field:", err)
		return
	}

	err = schema.AddField(*emailField)
	if err != nil {
		fmt.Println("Error adding email field:", err)
		return
	}

	// Try to add a duplicate field
	duplicateField := &Field{
		Name:     "name",
		Type:     types.TypeString,
		Required: true,
	}

	err = schema.AddField(*duplicateField)
	if err != nil {
		fmt.Println("Error adding duplicate field:", err)
	}

	schema.PrintSchema()
	storage.SaveStructToFile(schema, "./schemas/user.json")

	// Example valid document
	doc := map[string]interface{}{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}

	err = storage.LoadStructFromFile("./schemas/user.json", schema)
	if err != nil {
		fmt.Println("Error loading schema:", err)
		return
	}

	schema.PrintSchema()

	err = schema.AddRecord(doc)
	if err != nil {
		fmt.Println("Error adding record:", err)
		return
	}
	fmt.Println("Record added successfully")

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

// Build a schema
func Build() *Schema {
	schema := NewSchema()
	// schema.AddField("email", types.TypeString, false)
	return schema
}
