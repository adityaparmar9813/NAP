package driver

import (
	"fmt"

	"github.com/adityaparmar9813/NAP/internal/schema"
	"github.com/adityaparmar9813/NAP/internal/storage"
)

func Collection() (*schema.Schema, error) {
	schema := schema.NewSchema()

	err := storage.LoadStructFromFile("./schemas/user.json", schema)
	if err != nil {
		return schema, fmt.Errorf("error loading schema: %w", err)
	}

	return schema, nil
}
