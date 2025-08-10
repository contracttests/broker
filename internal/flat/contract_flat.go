package flat

import (
	"github.com/contracttests/broker/internal/dsl"
)

type FlatContract struct {
	Resources []FlatResource
	Schemas   FlatSchemas
}

func Contract(contractDsl dsl.Contract) FlatContract {
	flatResources := Resources(contractDsl)
	flatSchemas := Schemas(contractDsl)

	return FlatContract{
		Resources: flatResources,
		Schemas:   flatSchemas,
	}
}
