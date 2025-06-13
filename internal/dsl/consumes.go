package dsl

import (
	"fmt"
	"strconv"
	"strings"
)

type Consumes map[string]Consume

type Consume struct {
	Rest Resources `json:"rest,omitzero"`
}

type ComparableConsumer struct {
	Path         string
	SchemaName   string
	ServiceName  string
	ApiType      string
	ResourcePath string
	Method       string
	StatusCode   int
	ConsumerType string
}

func NewComparableConsumers(unknown Consumes) []ComparableConsumer {
	return recursiveComparableConsumersHelper([]ComparableConsumer{}, "consumes", unknown)
}

func newComparableConsumer(path string, schemaName string) ComparableConsumer {
	requestBodyParams := strings.Split(path, ";")

	if strings.Contains(path, "requestBody") {
		_, serviceName, apiType, resourcePath, method, _ :=
			requestBodyParams[0],
			requestBodyParams[1],
			requestBodyParams[2],
			requestBodyParams[3],
			requestBodyParams[4],
			requestBodyParams[5]

		return ComparableConsumer{
			Path:         path,
			SchemaName:   schemaName,
			ServiceName:  serviceName,
			ApiType:      apiType,
			ResourcePath: resourcePath,
			Method:       method,
			ConsumerType: "requestBody",
		}
	}

	_, serviceName, apiType, resourcePath, method, _, statusCodeAsString :=
		requestBodyParams[0],
		requestBodyParams[1],
		requestBodyParams[2],
		requestBodyParams[3],
		requestBodyParams[4],
		requestBodyParams[5],
		requestBodyParams[6]

	statusCode, err := strconv.Atoi(statusCodeAsString)
	if err != nil {
		panic(fmt.Sprintf("invalid status code: %s", statusCodeAsString))
	}

	return ComparableConsumer{
		Path:         path,
		SchemaName:   schemaName,
		ServiceName:  serviceName,
		ApiType:      apiType,
		ResourcePath: resourcePath,
		Method:       method,
		ConsumerType: "response",
		StatusCode:   statusCode,
	}
}

func recursiveComparableConsumersHelper(resources []ComparableConsumer, path string, unknown any) []ComparableConsumer {
	switch unknown := unknown.(type) {
	case Consumes:
		for name, resource := range unknown {
			consumesServicePathName := fmt.Sprintf("%s;%s", path, name)
			resources = recursiveComparableConsumersHelper(
				resources,
				consumesServicePathName,
				resource,
			)
		}
		return resources

	case Consume:
		for name, resource := range unknown.Rest {
			consumeRestPathName := fmt.Sprintf("%s;rest;%s", path, name)
			resources = recursiveComparableConsumersHelper(
				resources,
				consumeRestPathName,
				resource,
			)
		}

		return resources

	case Resource:
		if unknown.Get.IsNonZero() {
			getMethodPathName := fmt.Sprintf("%s;%s", path, "get")
			resources = recursiveComparableConsumersHelper(
				resources,
				getMethodPathName,
				unknown.Get,
			)
		}

		if unknown.Post.IsNonZero() {
			postMethodPathName := fmt.Sprintf("%s;%s", path, "post")
			resources = recursiveComparableConsumersHelper(
				resources,
				postMethodPathName,
				unknown.Post,
			)
		}

		if unknown.Put.IsNonZero() {
			putMethodPathName := fmt.Sprintf("%s;%s", path, "put")
			resources = recursiveComparableConsumersHelper(
				resources,
				putMethodPathName,
				unknown.Put,
			)
		}

		if unknown.Delete.IsNonZero() {
			deleteMethodPathName := fmt.Sprintf("%s;%s", path, "delete")
			resources = recursiveComparableConsumersHelper(
				resources,
				deleteMethodPathName,
				unknown.Delete,
			)
		}

		return resources

	case GetMethod:
		responsePathName := fmt.Sprintf("%s;%s", path, "response")
		return recursiveComparableConsumersHelper(
			resources,
			responsePathName,
			unknown.Responses,
		)

	case PostMethod:
		responsePathName := fmt.Sprintf("%s;%s", path, "response")
		resources = recursiveComparableConsumersHelper(
			resources,
			responsePathName,
			unknown.Responses,
		)
		if unknown.HasRequestBody() {
			requestBodyPathName := fmt.Sprintf("%s;%s", path, "requestBody")
			resources = append(
				resources,
				newComparableConsumer(requestBodyPathName, unknown.RequestBody),
			)
		}
		return resources

	case PutMethod:
		responsePathName := fmt.Sprintf("%s;%s", path, "response")
		resources = recursiveComparableConsumersHelper(
			resources,
			responsePathName,
			unknown.Responses,
		)
		if unknown.HasRequestBody() {
			requestBodyPathName := fmt.Sprintf("%s;%s", path, "requestBody")
			resources = append(
				resources,
				newComparableConsumer(requestBodyPathName, unknown.RequestBody),
			)
		}
		return resources

	case DeleteMethod:
		responsePathName := fmt.Sprintf("%s;%s", path, "response")
		return recursiveComparableConsumersHelper(
			resources,
			responsePathName,
			unknown.Responses,
		)

	case Responses:
		for code, schemaName := range unknown {
			statusCodePathName := fmt.Sprintf("%s;%d", path, code)
			resources = append(
				resources,
				newComparableConsumer(statusCodePathName, schemaName),
			)
		}
		return resources

	}

	return resources
}
