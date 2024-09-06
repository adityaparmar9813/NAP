package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adityaparmar9813/NAP/internal/schema"
	"github.com/adityaparmar9813/NAP/internal/storage"
	"github.com/adityaparmar9813/NAP/internal/types"
)

type Field = schema.Field

type MockValidator struct{}

func (mv MockValidator) ValidateType(value interface{}, fieldType types.FieldType) error {
	// Mock validation: assume all types are valid
	return nil
}

type MockStorage struct {
	Data map[string][]byte
}

func NewMockStorage() *MockStorage {
	return &MockStorage{Data: make(map[string][]byte)}
}

func (ms *MockStorage) SaveStructToFile(v interface{}, filename string) error {
	data, err := storage.StructToJSON(v)
	if err != nil {
		return err
	}
	ms.Data[filename] = data
	return nil
}

func (ms *MockStorage) LoadStructFromFile(filename string, v interface{}) error {
	data, exists := ms.Data[filename]
	if !exists {
		return os.ErrNotExist
	}
	return storage.JSONToStruct(data, v)
}

func TestAddField(t *testing.T) {
	schema := schema.NewSchema("test_schema")
	field := Field{Name: "name", Type: types.TypeString, Required: true}

	err := schema.AddField(field)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, exists := schema.Fields["name"]; !exists {
		t.Fatalf("expected field 'name' to exist in schema")
	}
}

func TestAddField_DuplicateField(t *testing.T) {
	schema := schema.NewSchema("test_schema")
	field := Field{Name: "name", Type: types.TypeString, Required: true}

	_ = schema.AddField(field)
	err := schema.AddField(field) // Adding the same field again

	if err == nil {
		t.Fatalf("expected error due to duplicate field, got none")
	}
}

func TestValidate(t *testing.T) {
	schema := schema.NewSchema("test_schema")
	schema.AddField(Field{Name: "name", Type: types.TypeString, Required: true})
	schema.AddField(Field{Name: "age", Type: types.TypeInt, Required: false})

	validator := MockValidator{}
	doc := map[string]interface{}{"name": "John Doe"}

	err := schema.Validate(doc, validator)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingRequiredField(t *testing.T) {
	schema := schema.NewSchema("test_schema")
	schema.AddField(Field{Name: "name", Type: types.TypeString, Required: true})

	validator := MockValidator{}
	doc := map[string]interface{}{"age": 30} // Missing required 'name' field

	err := schema.Validate(doc, validator)
	if err == nil {
		t.Fatalf("expected error due to missing required field, got none")
	}
}

func TestAddRecord(t *testing.T) {
	mockStorage := NewMockStorage()
	mockValidator := MockValidator{}
	schema, _ := schema.BuildSchema("test_schema", mockStorage)

	doc := map[string]interface{}{"name": "John Doe", "age": 30}

	err := schema.AddRecord(doc, mockValidator, mockStorage)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if the document was saved correctly
	collectionPath := filepath.Join("./collections", "test_schema")
	fileName := filepath.Join(collectionPath, doc["uuid"].(string)+".json")

	if _, exists := mockStorage.Data[fileName]; !exists {
		t.Fatalf("expected record to be saved, but it wasn't")
	}
}

func TestGetRecord(t *testing.T) {
	// mockStorage := NewMockStorage()
	// mockValidator := MockValidator{}
	// schema, _ := schema.BuildSchema("test_schema", mockStorage)

	// doc1 := map[string]interface{}{"name": "John Doe", "age": 30}
	// doc2 := map[string]interface{}{"name": "Jane Doe", "age": 25}

	// _ = schema.AddRecord(doc1, mockValidator, mockStorage)
	// _ = schema.AddRecord(doc2, mockValidator, mockStorage)

	// ageCriteria := map[string]interface{}{
	// 	"age": 25,
	// }
	// records, err := schema.GetRecord(ageCriteria, mockStorage)
	// if err != nil {
	// 	t.Fatalf("expected no error, got %v", err)
	// }

	// if len(records) != 1 || !reflect.DeepEqual(records[0]["name"], "Jane Doe") {
	// 	t.Fatalf("expected to get Jane Doe, got %v", records)
	// }
}

func TestGetRecord_NoMatchingRecords(t *testing.T) {
	mockStorage := NewMockStorage()
	mockValidator := MockValidator{}
	schema, _ := schema.BuildSchema("test_schema", mockStorage)

	doc := map[string]interface{}{"name": "John Doe", "age": 30}
	_ = schema.AddRecord(doc, mockValidator, mockStorage)

	criteria := map[string]interface{}{"age": 99}
	records, err := schema.GetRecord(criteria, mockStorage)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(records) != 0 {
		t.Fatalf("expected no matching records, got %v", records)
	}
}

func TestPrintSchema(t *testing.T) {
	schema := schema.NewSchema("test_schema")
	schema.AddField(Field{Name: "name", Type: types.TypeString, Required: true})
	schema.AddField(Field{Name: "age", Type: types.TypeInt, Required: false})

	schema.PrintSchema()
	// Since this prints to the console, you would manually verify the output,
	// or capture the output using a custom logger in real-world scenarios.
}
