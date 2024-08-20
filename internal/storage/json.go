package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// StructToJSON converts a struct to a JSON byte slice
func StructToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JSONToStruct converts a JSON byte slice back to a struct
func JSONToStruct(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// SaveJSONToFile saves JSON data to a file
func SaveJSONToFile(data []byte, filename string) error {
	// Ensure the directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create or open the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the data
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// LoadJSONFromFile loads JSON data from a file
func LoadJSONFromFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// SaveStructToFile converts a struct to JSON and saves it to a file
func SaveStructToFile(v interface{}, filename string) error {
	data, err := StructToJSON(v)
	if err != nil {
		return fmt.Errorf("failed to convert struct to JSON: %w", err)
	}

	return SaveJSONToFile(data, filename)
}

// LoadStructFromFile loads JSON data from a file and converts it to a struct
func LoadStructFromFile(filename string, v interface{}) error {
	data, err := LoadJSONFromFile(filename)
	if err != nil {
		return err
	}

	return JSONToStruct(data, v)
}
