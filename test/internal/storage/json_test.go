package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/adityaparmar9813/NAP/internal/storage"
)

type TestStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestTestSaveStructToFile(t *testing.T) {
	fs := storage.NewFileStorage()
	testData := TestStruct{Name: "Test", Value: 42}
	filename := filepath.Join(os.TempDir(), "test_save_struct.json")

	// Clean up after test
	defer os.Remove(filename)

	err := fs.SaveStructToFile(&testData, filename)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("expected file to exist, but it doesn't")
	}
}

func TestTestLoadStructFromFile(t *testing.T) {
	fs := storage.NewFileStorage()
	testData := TestStruct{Name: "Test", Value: 42}
	filename := filepath.Join(os.TempDir(), "test_load_struct.json")

	// Save the struct to the file
	err := fs.SaveStructToFile(&testData, filename)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Load the struct from the file
	var loadedData TestStruct
	err = fs.LoadStructFromFile(filename, &loadedData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if the loaded data is the same as the original
	if !reflect.DeepEqual(testData, loadedData) {
		t.Fatalf("expected %v, got %v", testData, loadedData)
	}
}

func TestTestSaveStructToFile_DirectoryCreation(t *testing.T) {
	fs := storage.NewFileStorage()
	testData := TestStruct{Name: "TestDir", Value: 123}
	dirname := filepath.Join(os.TempDir(), "nested/dir/structure")
	filename := filepath.Join(dirname, "test_save_struct.json")

	// Clean up after test
	defer os.RemoveAll(dirname)

	err := fs.SaveStructToFile(&testData, filename)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("expected file to exist, but it doesn't")
	}
}

func TestTestLoadStructFromFile_FileNotFound(t *testing.T) {
	fs := storage.NewFileStorage()
	filename := filepath.Join(os.TempDir(), "non_existent_file.json")

	var loadedData TestStruct
	err := fs.LoadStructFromFile(filename, &loadedData)

	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestTestSaveAndLoad_EmptyStruct(t *testing.T) {
	fs := storage.NewFileStorage()
	testData := TestStruct{}
	filename := filepath.Join(os.TempDir(), "test_empty_struct.json")

	// Clean up after test
	defer os.Remove(filename)

	// Save the empty struct
	err := fs.SaveStructToFile(&testData, filename)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Load the empty struct
	var loadedData TestStruct
	err = fs.LoadStructFromFile(filename, &loadedData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check if the loaded data is the same as the original
	if !reflect.DeepEqual(testData, loadedData) {
		t.Fatalf("expected %v, got %v", testData, loadedData)
	}
}

func TestTestSaveStructToFile_InvalidJSON(t *testing.T) {
	fs := storage.NewFileStorage()
	testData := make(chan int) // Channels can't be serialized to JSON
	filename := filepath.Join(os.TempDir(), "test_invalid_json.json")

	// Clean up after test
	defer os.Remove(filename)

	err := fs.SaveStructToFile(testData, filename)

	if err == nil {
		t.Fatalf("expected error due to invalid JSON, got none")
	}
}

func TestTestLoadStructFromFile_InvalidJSON(t *testing.T) {
	fs := storage.NewFileStorage()
	filename := filepath.Join(os.TempDir(), "test_invalid_json.json")

	// Write invalid JSON data to the file
	invalidJSON := []byte(`{"name": "Test", "value": }`)
	if err := os.WriteFile(filename, invalidJSON, 0644); err != nil {
		t.Fatalf("failed to write invalid JSON to file: %v", err)
	}

	// Clean up after test
	defer os.Remove(filename)

	var loadedData TestStruct
	err := fs.LoadStructFromFile(filename, &loadedData)

	if err == nil {
		t.Fatalf("expected error due to invalid JSON, got none")
	}
}
func TestFileStorage_LoadStructFromFile_InvalidJSON(t *testing.T) {
	fs := storage.NewFileStorage()
	filename := filepath.Join(os.TempDir(), "test_invalid_json.json")

	// Write invalid JSON data to the file
	invalidJSON := []byte(`{"name": "Test", "value": }`)
	if err := os.WriteFile(filename, invalidJSON, 0644); err != nil {
		t.Fatalf("failed to write invalid JSON to file: %v", err)
	}

	// Clean up after test
	defer os.Remove(filename)

	var loadedData TestStruct
	err := fs.LoadStructFromFile(filename, &loadedData)

	if err == nil {
		t.Fatalf("expected error due to invalid JSON, got none")
	}
}
