package dsl

type Contract struct {
	Api      Api      `json:"api,omitzero"`
	Provides Provides `json:"provides,omitzero"`
	Consumes Consumes `json:"consumes,omitzero"`
	Schemas  Schemas  `json:"schemas,omitzero"`
}

type Api struct {
	Name string `json:"name,omitzero"`
}

type ComparableContract struct {
	Schemas           ComparableSchemas
	ProvidesResources []ComparableProvider
	ConsumesResources []ComparableConsumer
}

func NewComparableContract(contract Contract) ComparableContract {
	return ComparableContract{
		Schemas:           NewComparableSchemas(contract.Schemas),
		ProvidesResources: NewComparableProviders(contract.Provides),
		ConsumesResources: NewComparableConsumers(contract.Consumes),
	}
}
