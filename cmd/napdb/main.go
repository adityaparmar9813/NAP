package main

import (
	"github.com/adityaparmar9813/NAP/internal/schema"
	"github.com/adityaparmar9813/NAP/internal/storage"
	"github.com/adityaparmar9813/NAP/internal/validator"
)

func main() {
	// Initialize the storage and validator implementations
	storageImpl := storage.NewFileStorage()   // Create an instance of FileStorage
	validatorImpl := validator.NewValidator() // Create an instance of Validator

	// Call the Test function with the required interfaces
	schema.Test(storageImpl, validatorImpl)

	// Example of other operations
	// schema, err := driver.Collection()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// schema.PrintSchema()
}
