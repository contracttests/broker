package flat

import "github.com/contracttests/broker/internal/model"

type FlatSchemas map[string]FlatSchema

type FlatSchema []FlatProperty

type FlatProperty struct {
	FullPath string
	Type     string
}

func (f *FlatProperty) IsPrimitive() bool {
	return f.Type == "string" ||
		f.Type == "number" ||
		f.Type == "boolean" ||
		f.Type == "integer" ||
		f.Type == "object" ||
		f.Type == "array"
}

func NewSchema(
	hash string,
	flatSchema FlatSchema,
) model.Schema {
	properties := make(map[string]model.Property)

	for _, property := range flatSchema {
		properties[property.FullPath] = model.NewProperty(property.FullPath, property.Type)
	}

	return model.Schema{
		Hash:       hash,
		Properties: properties,
	}
}
