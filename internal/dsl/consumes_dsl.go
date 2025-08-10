package dsl

type Consumes struct {
	Rest    Rest              `json:"rest,omitzero"`
	Message map[string]string `json:"message,omitzero"`
}

type ConsumesServices map[string]Consumes

// func SimplifyConsumers(serviceName string, unknown map[string]Consume) []model.Consumer {
// 	fullPath := fmt.Sprintf("%s;consumes", serviceName)
// 	return simplifyConsumersRecursively([]model.Consumer{}, fullPath, unknown)
// }

// func newSimplifiedConsumer(fullPath string, schemaName string) model.Consumer {
// 	parts := strings.Split(fullPath, ";")
