package types

// FieldType represents the type of a field
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeFloat   FieldType = "float"
	TypeBoolean FieldType = "boolean"
)

// ValidatorFunc is a function type for custom validators
type ValidatorFunc func(interface{}) bool
