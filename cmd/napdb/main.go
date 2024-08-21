package main

import (
	"fmt"

	"github.com/adityaparmar9813/NAP/internal/driver"
	"github.com/adityaparmar9813/NAP/internal/schema"
)

func main() {
	schema.Test()

	schema, err := driver.Collection()
	if err != nil {
		fmt.Println(err)
	}
	schema.PrintSchema()
}
