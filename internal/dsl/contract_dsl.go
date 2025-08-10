package dsl

type Contract struct {
	Api              Api              `json:"api,omitzero"`
	Provides         Provides         `json:"provides,omitzero"`
	ConsumesServices ConsumesServices `json:"consumes,omitzero"`
	Schemas          Schemas          `json:"schemas,omitzero"`
}

type Api struct {
	Name string `json:"name,omitzero"`
}

// func Simplify(contract Contract) {
// 	serviceName := contract.Api.Name

// 	// simplifiedRestRequestBodyProviders := SimplifyProviders(serviceName, contract.Provides)
// 	// simplifiedRestRequestBodyConsumers := SimplifyConsumers(serviceName, contract.Consumes)

// 	// simplifiedSchemas := []model.Schema{}

// 	// for _, simplifiedProvider := range simplifiedProviders {
// 	// 	simplifiedSchema := DslToSchema(
// 	// 		simplifiedProvider.SchemaName,
// 	// 		simplifiedProvider.FullPath,
// 	// 		contract.Schemas[simplifiedProvider.SchemaName],
// 	// 	)
// 	// 	simplifiedSchemas = append(simplifiedSchemas, simplifiedSchema)
// 	// }

// 	// for _, simplifiedConsumer := range simplifiedConsumers {
// 	// 	simplifiedSchema := DslToSchema(
// 	// 		simplifiedConsumer.SchemaName,
// 	// 		simplifiedConsumer.FullPath,
// 	// 		contract.Schemas[simplifiedConsumer.SchemaName],
// 	// 	)
// 	// 	simplifiedSchemas = append(simplifiedSchemas, simplifiedSchema)
// 	// }

// 	// return model.Contract{
// 	// 	// Schemas:   simplifiedSchemas,
// 	// 	RestRequestBodyProviders: simplifiedProviders,
// 	// 	RestRequestBodyConsumers: simplifiedConsumers,
// 	// }
// }
