package tool

import "github.com/contracttests/broker/internal/model"

func SchemaDiff(left model.Schema, right model.Schema) model.Schema {
	diff := model.Schema{
		Hash:       left.Hash,
		Properties: make(map[string]model.Property),
	}

	for propertyFullPath, property := range left.Properties {
		if _, ok := right.Properties[propertyFullPath]; !ok {
			diff.Properties[propertyFullPath] = property
		}
	}

	return diff
}
