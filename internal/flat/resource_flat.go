package flat

import (
	"strconv"
	"strings"

	"github.com/contracttests/broker/internal/dsl"
)

type FlatResource struct {
	FullPath   string
	SchemaName string
}

func newResourceFullPath(parts ...string) string {
	resourceParthSeparator := ";"

	return strings.Join(parts, resourceParthSeparator)
}

func Resources(contractDsl dsl.Contract) []FlatResource {
	return buildFlatResources([]FlatResource{}, "", contractDsl)
}

func buildFlatResources(flatResources []FlatResource, fullPath string, unknown any) []FlatResource {
	switch unknown := unknown.(type) {
	case dsl.Contract:
		fullPath = unknown.Api.Name

		for serviceName, consumes := range unknown.ConsumesServices {
			newFullPath := newResourceFullPath(fullPath, "consumes", serviceName)
			flatResources = buildFlatResources(flatResources, newFullPath, consumes)
		}

		newFullPath := newResourceFullPath(fullPath, "provides")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Provides)

		return flatResources

	case dsl.Consumes:
		flatResources = buildFlatResources(flatResources, fullPath, unknown.Rest)
		flatResources = buildFlatResources(flatResources, fullPath, unknown.Message)

		return flatResources

	case dsl.Provides:
		flatResources = buildFlatResources(flatResources, fullPath, unknown.Rest)
		flatResources = buildFlatResources(flatResources, fullPath, unknown.Message)

		return flatResources

	case dsl.Message:
		for messageName, schemaName := range unknown {
			newFullPath := newResourceFullPath(fullPath, "message", messageName)
			flatResources = append(flatResources, FlatResource{
				FullPath:   newFullPath,
				SchemaName: schemaName,
			})
		}

		return flatResources

	case dsl.Rest:
		for endpoint, methods := range unknown {
			if methods.Get.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildFlatResources(flatResources, newFullPath, methods.Get)
			}

			if methods.Post.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildFlatResources(flatResources, newFullPath, methods.Post)
			}

			if methods.Put.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildFlatResources(flatResources, newFullPath, methods.Put)
			}

			if methods.Delete.IsNonZero() {
				newFullPath := newResourceFullPath(fullPath, "rest", endpoint)
				flatResources = buildFlatResources(flatResources, newFullPath, methods.Delete)
			}
		}

	case dsl.GetMethod:
		newFullPath := newResourceFullPath(fullPath, "get", "responses")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case dsl.PostMethod:
		if unknown.HasRequestBody() {
			newFullPath := newResourceFullPath(fullPath, "post", "requestBody")
			flatResources = append(flatResources, FlatResource{
				FullPath:   newFullPath,
				SchemaName: unknown.RequestBody,
			})
		}

		newFullPath := newResourceFullPath(fullPath, "post", "responses")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case dsl.PutMethod:
		if unknown.HasRequestBody() {
			newFullPath := newResourceFullPath(fullPath, "put", "requestBody")
			flatResources = append(flatResources, FlatResource{
				FullPath:   newFullPath,
				SchemaName: unknown.RequestBody,
			})
		}

		newFullPath := newResourceFullPath(fullPath, "put", "responses")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case dsl.DeleteMethod:
		newFullPath := newResourceFullPath(fullPath, "delete", "responses")
		flatResources = buildFlatResources(flatResources, newFullPath, unknown.Responses)

		return flatResources

	case dsl.Responses:
		for statusCode, schemaName := range unknown {
			newFullPath := newResourceFullPath(fullPath, strconv.Itoa(statusCode))
			flatResources = append(flatResources, FlatResource{
				FullPath:   newFullPath,
				SchemaName: schemaName,
			})
		}

		return flatResources
	}

	return flatResources
}
