package flat

type FlatContract struct {
	Resources []FlatResource
	Schemas   FlatSchemas
}

// func (flatContract FlatContract) ToModelContract() model.Contract {
// 	contract := model.Contract{
// 		Resources: []model.Resource{},
// 		Schemas:   make(map[string]model.Schema),
// 	}

// 	for _, flatResource := range flatContract.Resources {
// 		resource := NewResource(flatResource)
// 		contract.Resources = append(contract.Resources, resource)

// 		flatSchema := flatContract.Schemas[flatResource.SchemaName]

// 		schema := NewSchema(resource.UniqueHash, flatSchema)
// 		contract.Schemas[resource.UniqueHash] = schema
// 	}

// 	return contract
// }
