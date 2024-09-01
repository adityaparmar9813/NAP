package schema

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adityaparmar9813/NAP/internal/storage"
	"github.com/adityaparmar9813/NAP/internal/types"
	"github.com/adityaparmar9813/NAP/internal/validator"
	"github.com/google/uuid"
)

type Field struct {
	Name     string
	Type     types.FieldType
	Required bool
}

type Schema struct {
	Name   string
	Fields map[string]Field
}

func NewSchema(name string) *Schema {
	return &Schema{
		Name:   name,
		Fields: make(map[string]Field),
	}
}

type SchemaInterface interface {
	AddField(field Field) error
	Validate(doc map[string]interface{}) error
	PrintSchema()
	AddRecord(doc map[string]interface{}, storage storage.StorageInterface) error
	GetRecord(uuid string, storage storage.StorageInterface) (map[string]interface{}, error)
}

func (s *Schema) AddField(field Field) error {
	if _, exists := s.Fields[field.Name]; exists {
		return fmt.Errorf("field '%s' already exists in schema", field.Name)
	}

	s.Fields[field.Name] = field

	return nil
}

func BuildSchema(name string, storage storage.StorageInterface, fields ...Field) (*Schema, error) {
	schema := NewSchema(name)

	// Add UUID field by default
	uuidField := Field{
		Name:     "uuid",
		Type:     types.TypeString,
		Required: true,
	}
	err := schema.AddField(uuidField)
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		err := schema.AddField(field)
		if err != nil {
			return nil, err
		}
	}

	// Save the schema to a file
	schemaPath := filepath.Join("./schemas", name+".json")
	err = storage.SaveStructToFile(schema, schemaPath)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func (s *Schema) Validate(doc map[string]interface{}, validator validator.ValidatorInterface) error {
	for fieldName, field := range s.Fields {
		// Skip validation for uuid field
		if fieldName == "uuid" {
			continue
		}

		value, exists := doc[fieldName]

		if !exists {
			if field.Required {
				return fmt.Errorf("required field '%s' is missing", fieldName)
			}
			continue
		}

		if err := validator.ValidateType(value, field.Type); err != nil {
			return fmt.Errorf("field '%s': %w", fieldName, err)
		}
	}

	return nil
}

func (s *Schema) AddRecord(doc map[string]interface{}, validator validator.ValidatorInterface, storage storage.StorageInterface) error {
	// Validate the document before adding the UUID
	err := s.Validate(doc, validator)
	if err != nil {
		return err
	}

	// Add UUID field
	recordID := uuid.New().String()
	doc["uuid"] = recordID

	// Create collection directory if it doesn't exist
	collectionPath := filepath.Join("./collections", s.Name)
	err = os.MkdirAll(collectionPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create collection directory: %w", err)
	}

	// Save the record to a file named after its UUID
	filePath := filepath.Join(collectionPath, recordID+".json")
	err = storage.SaveStructToFile(doc, filePath)
	if err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}

	return nil
}

func (s *Schema) GetRecord(criteria map[string]interface{}, storage storage.StorageInterface) ([]map[string]interface{}, error) {
	collectionPath := filepath.Join("./collections", s.Name)
	files, err := os.ReadDir(collectionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection directory: %w", err)
	}

	var matchingRecords []map[string]interface{}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(collectionPath, file.Name())
		var record map[string]interface{}
		err := storage.LoadStructFromFile(filePath, &record)
		if err != nil {
			return nil, fmt.Errorf("failed to load record from file %s: %w", file.Name(), err)
		}

		if validator.MatchesCriteria(record, criteria) {
			matchingRecords = append(matchingRecords, record)
		}
	}

	return matchingRecords, nil
}

func (s *Schema) PrintSchema() {
	fmt.Printf("Schema for collection '%s':\n", s.Name)
	for name, field := range s.Fields {
		fmt.Printf("%s (%s, required=%t)\n", name, field.Type, field.Required)
	}
}

func Test(storage storage.StorageInterface, validator validator.ValidatorInterface) {
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

	schema, err := BuildSchema("users", storage, nameField, ageField, emailField)
	if err != nil {
		fmt.Println("Error building schema:", err)
		return
	}

	schema.PrintSchema()

	docs := []map[string]interface{}{
		{"name": "Arpit Dubey", "age": 22, "email": "adubey_be21@thapar.edu"},
		{"name": "Ansh Bajaj", "age": 20, "email": "anshbajaj07@gmail.com"},
	}

	for _, doc := range docs {
		err = schema.AddRecord(doc, validator, storage)
		if err != nil {
			fmt.Println("Error adding record:", err)
			return
		}
		fmt.Println("Record added successfully")

	}

	ageCriteria := map[string]interface{}{
		"age": 20,
	}
	matchingRecords, err := schema.GetRecord(ageCriteria, storage)
	if err != nil {
		fmt.Println("Error getting records:", err)
		return
	}
	fmt.Println("Matching records:", matchingRecords)

	// Invalid document examples
	invalidDocs := []map[string]interface{}{
		{"name": "Bob", "email": "@invalid"}, // Required field age is missing
		{"age": 25},                          // Missing required field (name)
	}

	for i, doc := range invalidDocs {
		err = schema.Validate(doc, validator)
		if err != nil {
			fmt.Printf("Invalid document %d: %v\n", i+1, err)
		}
	}
}
