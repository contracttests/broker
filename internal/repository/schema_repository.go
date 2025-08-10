package repository

import "github.com/contracttests/broker/internal/model"

var schemasMap = make(map[string]model.Schema)

func SaveSchema(schema model.Schema) {
	schemasMap[schema.Hash] = schema
}

func GetSchema(hash string) model.Schema {
	var schema model.Schema
	if schema, ok := schemasMap[hash]; ok {
		return schema
	}

	return schema
}
