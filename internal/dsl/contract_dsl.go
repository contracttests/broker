package dsl

import (
	"fmt"
	"strconv"

	"github.com/contracttests/broker/server/internal/model"
)

type Contract struct {
	Name             string              `json:"name,omitzero"`
	Owner            string              `json:"owner"`
	Provides         Provides            `json:"provides,omitzero"`
	ConsumesServices ConsumesServicesMap `json:"consumes,omitzero"`
	Schemas          SchemasMap          `json:"schemas,omitzero"`
}

func (c *Contract) ToContractModel() model.Contract {
	contract := model.Contract{
		Name:  c.Name,
		Owner: c.Owner,
	}

	c.buildResources(&contract, NewResourcePath(""), *c)
	return contract
}

func (c *Contract) buildResources(contractModel *model.Contract, resourcePath ResourcePath, unknown any) {
	switch unknown := unknown.(type) {
	case Contract:
		dsl := unknown

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

			properties := buildSchema(
				NewDepthCounter(postMethod.RequestBody),
				c.Schemas,
				make(map[string]model.Property),
				NewPropertyPath("root"),
				c.Schemas[postMethod.RequestBody],
			)

			contractModel.AddResource(requestResourcePath.ToResource(properties))
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

			properties := buildSchema(
				NewDepthCounter(putMethod.RequestBody),
				c.Schemas,
				make(map[string]model.Property),
				NewPropertyPath("root"),
				c.Schemas[putMethod.RequestBody],
			)

			contractModel.AddResource(requestResourcePath.ToResource(properties))
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
			properties := buildSchema(
				NewDepthCounter(schemaName),
				c.Schemas,
				make(map[string]model.Property),
				NewPropertyPath("root"),
				c.Schemas[schemaName],
			)

			contractModel.AddResource(responseResourcePath.ToResource(properties))
		}
	}
}

func buildSchema(
	dethCounter *DepthCounter,
	schemas SchemasMap,
	properties map[string]model.Property,
	propertyPath PropertyPath,
	unknown any,
) map[string]model.Property {
	switch unknown := unknown.(type) {
	case Schema:
		if unknown.IsObject() {
			properties[propertyPath.String()] = model.NewProperty(propertyPath.String(), "object", unknown.Optional)

			for name, schemaProperties := range unknown.Properties {
				dethCounter.Enter()
				properties = buildSchema(
					dethCounter,
					schemas,
					properties,
					propertyPath.Append(name),
					schemaProperties,
				)
			}

			return properties
		}

		if unknown.IsArray() {
			properties[propertyPath.String()] = model.NewProperty(propertyPath.String(), "array", unknown.Optional)

			dethCounter.Enter()
			properties = buildSchema(
				dethCounter,
				schemas,
				properties,
				propertyPath.AppendArray(),
				unknown.Items,
			)

			return properties
		}

		if unknown.IsPrimitive() {
			properties[propertyPath.String()] = model.NewProperty(propertyPath.String(), unknown.Type, unknown.Optional)

			return properties
		}

		if unknown.IsRef() {
			dethCounter.Enter()
			properties = buildSchema(
				dethCounter,
				schemas,
				properties,
				propertyPath,
				schemas[unknown.Ref],
			)

			return properties
		}

		return properties
	case *Schema:
		dethCounter.Enter()
		return buildSchema(
			dethCounter,
			schemas,
			properties,
			propertyPath,
			*unknown,
		)
	default:
		panic(fmt.Sprintf("unknown schema type %T", unknown))
	}
}
