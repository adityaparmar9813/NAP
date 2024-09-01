package driver

// func Collection(storage schema.StorageInterface) (*schema.Schema, error) {
// 	schema := schema.NewSchema()

// 	err := storage.LoadStructFromFile("./schemas/user.json", schema)
// 	if err != nil {
// 		return schema, fmt.Errorf("error loading schema: %w", err)
// 	}

// 	return schema, nil
// }

// func Insert(storage schema.storageinterface, schema *schema.schema, data interface{}) error {
// 	err := schema.addrecord(data)
// 	if err != nil {
// 		return fmt.errorf("error inserting data: %w", err)
// 	}

// 	return nil
// }
