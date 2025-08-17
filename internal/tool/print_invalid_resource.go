package tool

import (
	"fmt"

	"github.com/contracttests/broker/internal/model"
)

func PrintInvalidConsumerResource(consumerResource model.Resource, providerResource model.Resource, schemaDiff model.Schema) {
	fmt.Println("--------------------------------")
	fmt.Printf("%s cannot provide at \n", providerResource.Provider.Name)
	fmt.Printf("Endpoint: %s\n", consumerResource.RestResource.Endpoint)
	fmt.Printf("Method: %s\n", consumerResource.RestResource.Method)
	if consumerResource.IsRequestBody() {
		fmt.Println("Request body: ")
	} else {
		fmt.Printf("Status Code: %s\n", consumerResource.RestResource.StatusCode)
		fmt.Println("Response:")
	}

	for _, property := range schemaDiff.Properties {
		fmt.Printf("  - The property %s is missing in the provider schema\n", property.Path)
	}
}

func PrintInvalidProviderResource(providerResource model.Resource, consumerResource model.Resource, schemaDiff model.Schema) {
	fmt.Println("--------------------------------")
	fmt.Printf("%s cannot be consume %s \n", consumerResource.Consumer.Name, providerResource.Provider.Name)
	fmt.Printf("Endpoint: %s\n", providerResource.RestResource.Endpoint)
	fmt.Printf("Method: %s\n", providerResource.RestResource.Method)
	if providerResource.IsRequestBody() {
		fmt.Println("Request body: ")
	} else {
		fmt.Printf("Status Code: %s\n", providerResource.RestResource.StatusCode)
		fmt.Println("Response:")
	}

	for _, property := range schemaDiff.Properties {
		fmt.Printf("  - The property %s is missing in the provider schema\n", property.Path)
	}
}
