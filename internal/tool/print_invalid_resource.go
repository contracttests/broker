package tool

import (
	"fmt"

	"github.com/contracttests/broker/internal/model"
)

func PrintInvalidConsumerResource(consumerRestResource model.Resource, providerRestResource model.Resource, schemaDiff model.Schema) {
	fmt.Println("--------------------------------")
	fmt.Printf("%s cannot provide at \n", consumerRestResource.ProviderName)
	fmt.Printf("Endpoint: %s\n", consumerRestResource.RestResource.Endpoint)
	fmt.Printf("Method: %s\n", consumerRestResource.RestResource.Method)
	if consumerRestResource.IsRequestBody() {
		fmt.Println("Request body: ")
	} else {
		fmt.Printf("Status Code: %s\n", consumerRestResource.RestResource.StatusCode)
		fmt.Println("Response:")
	}

	for _, property := range schemaDiff.Properties {
		fmt.Printf("  - The property %s is missing in the provider schema\n", property.Path)
	}
}

func PrintInvalidProviderResource(providerRestResource model.Resource, consumerRestResource model.Resource, schemaDiff model.Schema) {
	fmt.Println("--------------------------------")
	fmt.Printf("%s cannot be consume by %s \n", providerRestResource.ProviderName, consumerRestResource.ConsumerName)
	fmt.Printf("Endpoint: %s\n", providerRestResource.RestResource.Endpoint)
	fmt.Printf("Method: %s\n", providerRestResource.RestResource.Method)
	if providerRestResource.IsRequestBody() {
		fmt.Println("Request body: ")
	} else {
		fmt.Printf("Status Code: %s\n", providerRestResource.RestResource.StatusCode)
		fmt.Println("Response:")
	}

	for _, property := range schemaDiff.Properties {
		fmt.Printf("  - The property %s is missing in the provider schema\n", property.Path)
	}
}
