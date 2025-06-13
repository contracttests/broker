package dsl

import (
	"fmt"
	"strconv"
	"strings"
)

type Provides struct {
	Rest Resources `json:"rest,omitzero"`
}

type ComparableProvider struct {
	Path         string
	SchemaName   string
	ApiType      string
	ResourcePath string
	Method       string
	StatusCode   int
	ConsumerType string
}

func newComparableProvider(path string, schemaName string) ComparableProvider {
	requestBodyParams := strings.Split(path, ";")

	if strings.Contains(path, "requestBody") {
		_, apiType, resourcePath, method, _ :=
			requestBodyParams[0],
			requestBodyParams[1],
			requestBodyParams[2],
			requestBodyParams[3],
			requestBodyParams[4]

		return ComparableProvider{
			Path:         path,
			SchemaName:   schemaName,
			ApiType:      apiType,
			ResourcePath: resourcePath,
			Method:       method,
			ConsumerType: "requestBody",
		}
	}

	_, apiType, resourcePath, method, _, statusCodeAsString :=
		requestBodyParams[0],
		requestBodyParams[1],
		requestBodyParams[2],
		requestBodyParams[3],
		requestBodyParams[4],
		requestBodyParams[5]

	statusCode, err := strconv.Atoi(statusCodeAsString)
	if err != nil {
		panic(fmt.Sprintf("invalid status code: %s", statusCodeAsString))
	}

	return ComparableProvider{
		Path:         path,
		SchemaName:   schemaName,
		ApiType:      apiType,
		ResourcePath: resourcePath,
		Method:       method,
		StatusCode:   statusCode,
		ConsumerType: "response",
	}
}

func NewComparableProviders(provides Provides) []ComparableProvider {
	return recursiveComparableProvidersHelper([]ComparableProvider{}, "provides", provides)
}

func recursiveComparableProvidersHelper(resources []ComparableProvider, path string, unknown any) []ComparableProvider {
	switch unknown := unknown.(type) {
	case Provides:
		for name, resource := range unknown.Rest {
			restResourcePathName := fmt.Sprintf("%s;rest;%s", path, name)
			resources = recursiveComparableProvidersHelper(resources, restResourcePathName, resource)
		}
		return resources

	case Resource:
		if unknown.Get.IsNonZero() {
			getMethodPathName := fmt.Sprintf("%s;%s", path, "get")
			resources = recursiveComparableProvidersHelper(
				resources,
				getMethodPathName,
				unknown.Get,
			)
		}

		if unknown.Post.IsNonZero() {
			postMethodPathName := fmt.Sprintf("%s;%s", path, "post")
			resources = recursiveComparableProvidersHelper(
				resources,
				postMethodPathName,
				unknown.Post,
			)
		}

		if unknown.Put.IsNonZero() {
			putMethodPathName := fmt.Sprintf("%s;%s", path, "put")
			resources = recursiveComparableProvidersHelper(
				resources,
				putMethodPathName,
				unknown.Put,
			)
		}

		if unknown.Delete.IsNonZero() {
			deleteMethodPathName := fmt.Sprintf("%s;%s", path, "delete")
			resources = recursiveComparableProvidersHelper(
				resources,
				deleteMethodPathName,
				unknown.Delete,
			)
		}

		return resources

	case GetMethod:
		responsePathName := fmt.Sprintf("%s;%s", path, "response")
		return recursiveComparableProvidersHelper(
			resources,
			responsePathName,
			unknown.Responses,
		)

	case PostMethod:
		responsePathName := fmt.Sprintf("%s;%s", path, "response")
		resources = recursiveComparableProvidersHelper(
			resources,
			responsePathName,
			unknown.Responses,
		)
		if unknown.HasRequestBody() {
			requestBodyPathName := fmt.Sprintf("%s;%s", path, "requestBody")
			resources = append(
				resources,
				newComparableProvider(requestBodyPathName, unknown.RequestBody),
			)
		}
		return resources

	case PutMethod:
		resources = recursiveComparableProvidersHelper(resources, fmt.Sprintf("%s;%s", path, "response"), unknown.Responses)
		if unknown.HasRequestBody() {
			requestBodyPathName := fmt.Sprintf("%s;%s", path, "requestBody")
			resources = append(
				resources,
				newComparableProvider(requestBodyPathName, unknown.RequestBody),
			)
		}
		return resources

	case DeleteMethod:
		responsePathName := fmt.Sprintf("%s;%s", path, "response")
		return recursiveComparableProvidersHelper(
			resources,
			responsePathName,
			unknown.Responses,
		)

	case Responses:
		for code, schemaName := range unknown {
			statusCodePathName := fmt.Sprintf("%s;%d", path, code)
			resources = append(resources, newComparableProvider(statusCodePathName, schemaName))
		}
		return resources

	}

	return resources
}
