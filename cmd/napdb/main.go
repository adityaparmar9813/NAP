package main

import (
	"fmt"

	"github.com/adityaparmar9813/NAP/internal/schema"
	"github.com/adityaparmar9813/NAP/internal/types"
)

func main() {
	// Create a new schema
	schema := schema.NewSchema()

	// Add fields to the schema
	err := schema.AddField("name", types.TypeString, true, func(v interface{}) bool {
		s, ok := v.(string)
		return ok && len(s) > 0
	})
	if err != nil {
		fmt.Println("Error adding field:", err)
		return
	}

	err = schema.AddField("age", types.TypeInt, true, func(v interface{}) bool {
		age, ok := v.(int)
		return ok && age >= 0 && age <= 120
	})
	if err != nil {
		fmt.Println("Error adding field:", err)
		return
	}

	err = schema.AddField("email", types.TypeString, false, func(v interface{}) bool {
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
