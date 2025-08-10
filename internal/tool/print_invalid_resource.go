package tool

import (
	"fmt"

	"github.com/contracttests/broker/internal/model"
)

func PrintInvalidConsumerResource(consumerRestResource model.RestResource, providerRestResource model.RestResource, schemaDiff model.Schema) {
	fmt.Println("--------------------------------")
	fmt.Printf("%s cannot provide at \n", consumerRestResource.ProviderName)
	fmt.Printf("Endpoint: %s\n", consumerRestResource.Endpoint)
	fmt.Printf("Method: %s\n", consumerRestResource.Method)
	if consumerRestResource.IsRequestBody() {
		fmt.Println("Request body: ")
	} else {
		fmt.Printf("Status Code: %s\n", consumerRestResource.StatusCode)
		fmt.Println("Response:")
	}

	for _, property := range schemaDiff.Properties {
		fmt.Printf("  - The property %s is missing in the provider schema\n", property.Path)
	}
}

func PrintInvalidProviderResource(providerRestResource model.RestResource, consumerRestResource model.RestResource, schemaDiff model.Schema) {
	fmt.Println("--------------------------------")
	fmt.Printf("%s cannot be consume by %s \n", providerRestResource.ProviderName, consumerRestResource.ConsumerName)
	fmt.Printf("Endpoint: %s\n", providerRestResource.Endpoint)
	fmt.Printf("Method: %s\n", providerRestResource.Method)
	if providerRestResource.IsRequestBody() {
		fmt.Println("Request body: ")
	} else {
		fmt.Printf("Status Code: %s\n", providerRestResource.StatusCode)
		fmt.Println("Response:")
	}

	for _, property := range schemaDiff.Properties {
		fmt.Printf("  - The property %s is missing in the provider schema\n", property.Path)
	}
}
