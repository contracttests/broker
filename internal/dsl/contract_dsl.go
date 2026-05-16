package dsl

import (
	"fmt"
	"strconv"

	"github.com/contracttests/broker/server/internal/model"
)

type Api struct {
	Name string `json:"name,omitzero"`
}

type Contract struct {
	Api              Api                 `json:"api,omitzero"`
	Provides         Provides            `json:"provides,omitzero"`
	ConsumesServices ConsumesServicesMap `json:"consumes,omitzero"`
	Schemas          SchemasMap          `json:"schemas,omitzero"`
}

func (c *Contract) ToContractModel() model.Contract {
	contract := model.Contract{}
	c.buildResources(&contract, NewResourcePath(""), *c)
	return contract
}

func (c *Contract) buildResources(contractModel *model.Contract, resourcePath ResourcePath, unknown any) {
	switch unknown := unknown.(type) {
	case Contract:
		dsl := unknown
		resourcePath = resourcePath.Append(dsl.Api.Name)

		for serviceName, consumes := range dsl.ConsumesServices {
			consumerResourcePath := resourcePath.Append("consumes", serviceName)
			c.buildResources(contractModel, consumerResourcePath, consumes)
		}

		c.buildResources(contractModel, resourcePath.Append("provides"), dsl.Provides)

	case Consumes:
		consumes := unknown
		c.buildResources(contractModel, resourcePath, consumes.Rest)

	case Provides:
		provides := unknown
		c.buildResources(contractModel, resourcePath, provides.Rest)
		c.buildResources(contractModel, resourcePath, provides.Message)

	case Rest:
		rest := unknown
		for endpoint, methods := range rest {
			if methods.Get.IsNonZero() {
				c.buildResources(
					contractModel, 
					resourcePath.Append("rest", endpoint),
					methods.Get,
				)
			}

			if methods.Post.IsNonZero() {
				c.buildResources(
					contractModel, 
					resourcePath.Append("rest", endpoint),
					methods.Post,
				)
			}

			if methods.Put.IsNonZero() {
				c.buildResources(
					contractModel, 
					resourcePath.Append("rest", endpoint),
					methods.Put,
				)
			}

			if methods.Delete.IsNonZero() {
				c.buildResources(
					contractModel, 
					resourcePath.Append("rest", endpoint),
					methods.Delete,
				)
			}
		}

	case GetMethod:
		getMethod := unknown
		c.buildResources(
			contractModel, 
			resourcePath.Append("get", "responses"),
			getMethod.Responses,
		)

	case PostMethod:
		postMethod := unknown
		if postMethod.HasRequestBody() {
			requestResourcePath := resourcePath.Append("post", "request")

			schema := buildSchema(
				NewDepthCounter(postMethod.RequestBody),
				c.Schemas,
				model.NewSchema(),
				NewPropertyPath("root"),
				c.Schemas[postMethod.RequestBody],
			)

			if requestResourcePath.IsConsumer() {
				consumerRestRequestArgs := requestResourcePath.ToConsumerRestRequestArgs()
				consumerRestRequest := model.NewConsumerRestRequest(consumerRestRequestArgs, schema)
				contractModel.ConsumerRequests = append(contractModel.ConsumerRequests, consumerRestRequest)
			}

			if requestResourcePath.IsProvider() {
				providerRestRequestArgs := requestResourcePath.ToProviderRestRequestArgs()
				providerRestRequest := model.NewProviderRestRequest(providerRestRequestArgs, schema)
				contractModel.ProviderRequests = append(contractModel.ProviderRequests, providerRestRequest)
			}
		}

		c.buildResources(
			contractModel, 
			resourcePath.Append("post", "responses"),
			postMethod.Responses,
		)

	case PutMethod:
		putMethod := unknown
		if putMethod.HasRequestBody() {
			requestResourcePath := resourcePath.Append("put", "request")

			schema := buildSchema(
				NewDepthCounter(putMethod.RequestBody),
				c.Schemas,
				model.NewSchema(),
				NewPropertyPath("root"),
				c.Schemas[putMethod.RequestBody],
			)

			if requestResourcePath.IsConsumer() {
				consumerRestRequestArgs := requestResourcePath.ToConsumerRestRequestArgs()
				consumerRestRequest := model.NewConsumerRestRequest(consumerRestRequestArgs, schema)
				contractModel.ConsumerRequests = append(contractModel.ConsumerRequests, consumerRestRequest)
			}

			if requestResourcePath.IsProvider() {
				providerRestRequestArgs := requestResourcePath.ToProviderRestRequestArgs()
				providerRestRequest := model.NewProviderRestRequest(providerRestRequestArgs, schema)
				contractModel.ProviderRequests = append(contractModel.ProviderRequests, providerRestRequest)
			}
		}

		c.buildResources(
			contractModel, 
			resourcePath.Append("put", "responses"),
			putMethod.Responses,
		)

	case DeleteMethod:
		deleteMethod := unknown
		path := resourcePath.Append("delete", "responses")
		c.buildResources(contractModel, path, deleteMethod.Responses)

	case Responses:
		responses := unknown
		for statusCode, schemaName := range responses {
			responseResourcePath := resourcePath.Append(strconv.Itoa(statusCode))
			schema := buildSchema(
				NewDepthCounter(schemaName),
				c.Schemas,
				model.NewSchema(),
				NewPropertyPath("root"),
				c.Schemas[schemaName],
			)

			if responseResourcePath.IsConsumer() {
				consumerRestResponse := model.NewConsumerRestResponse(responseResourcePath.ToConsumerRestResponseArgs(), schema)
				contractModel.ConsumerResponses = append(contractModel.ConsumerResponses, consumerRestResponse)
			}

			if responseResourcePath.IsProvider() {
				providerRestResponse := model.NewProviderRestResponse(responseResourcePath.ToProviderRestResponseArgs(), schema)
				contractModel.ProviderResponses = append(contractModel.ProviderResponses, providerRestResponse)
			}
		}
	}
}

func buildSchema(
	dethCounter *DepthCounter,
	schemas SchemasMap,
	schema model.Schema,
	propertyPath PropertyPath,
	unknown any,
) model.Schema {
	switch unknown := unknown.(type) {
	case Schema:
		if unknown.IsObject() {
			schema.AddProperty(model.NewProperty(propertyPath.String(), "object", unknown.Optional))

			for name, schemaProperties := range unknown.Properties {
				dethCounter.Enter()
				schema = buildSchema(
					dethCounter,
					schemas,
					schema,
					propertyPath.Append(name),
					schemaProperties,
				)
			}

			return schema
		}

		if unknown.IsArray() {
			schema.AddProperty(model.NewProperty(propertyPath.String(), "array", unknown.Optional))

			dethCounter.Enter()
			schema = buildSchema(
				dethCounter,
				schemas,
				schema,
				propertyPath.AppendArray(),
				unknown.Items,
			)

			return schema
		}

		if unknown.IsPrimitive() {
			schema.AddProperty(model.NewProperty(propertyPath.String(), unknown.Type, unknown.Optional))

			return schema
		}

		if unknown.IsRef() {
			dethCounter.Enter()
			schema = buildSchema(
				dethCounter,
				schemas,
				schema,
				propertyPath,
				schemas[unknown.Ref],
			)

			return schema
		}

		return schema
	case *Schema:
		dethCounter.Enter()
		return buildSchema(
			dethCounter,
			schemas,
			schema,
			propertyPath,
			*unknown,
		)
	default:
		panic(fmt.Sprintf("unknown schema type %T", unknown))
	}
}
