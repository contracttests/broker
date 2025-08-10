package dsl

type Provides struct {
	Rest    Rest              `json:"rest,omitzero"`
	Message map[string]string `json:"message,omitzero"`
}

// func SimplifyProviders(serviceName string, provides Provides) []model.Provider {
// 	fullPath := fmt.Sprintf("%s;provides", serviceName)
// 	return simplifyProvidersRecursively([]model.Provider{}, fullPath, provides)
// }

// func newSimplifiedProvider(fullPath string, schemaName string) model.Provider {
// 	parts := strings.Split(fullPath, ";")

// 	if strings.Contains(fullPath, "requestBody") {
// 		providerServiceName, apiType, resourcePath, method :=
// 			parts[0],
// 			parts[2],
// 			parts[3],
// 			parts[4]

// 		return model.NewRequestBodyProvider(
// 			providerServiceName,
// 			apiType,
// 			resourcePath,
// 			method,
// 		)
// 	}

// 	providerServiceName, resourcePath, method, statusCodeAsString :=
// 		parts[0],
// 		parts[2],
// 		parts[3],
// 		parts[4],
// 		parts[6]

// 	statusCode, err := strconv.Atoi(statusCodeAsString)
// 	if err != nil {
// 		panic(fmt.Sprintf("invalid status code: %s", statusCodeAsString))
// 	}

// 	return model.NewRestResponseProvider(
// 		providerServiceName,
// 		resourcePath,
// 		method,
// 		statusCode,
// 		schemaName,
// 	)
// }

// func simplifyProvidersRecursively(simplifiedProviders []model.Provider, fullPath string, unknown any) []model.Provider {
// 	switch unknown := unknown.(type) {
// 	case Provides:
// 		for name, restResource := range unknown.Rest {
// 			newFullPath := fmt.Sprintf("%s;rest;%s", fullPath, name)
// 			simplifiedProviders = simplifyProvidersRecursively(simplifiedProviders, newFullPath, restResource)
// 		}
// 		return simplifiedProviders

// 	case Endpoint:
// 		if unknown.Get.IsNonZero() {
// 			newFullPath := fmt.Sprintf("%s;%s", fullPath, "get")
// 			simplifiedProviders = simplifyProvidersRecursively(
// 				simplifiedProviders,
// 				newFullPath,
// 				unknown.Get,
// 			)
// 		}

// 		if unknown.Post.IsNonZero() {
// 			newFullPath := fmt.Sprintf("%s;%s", fullPath, "post")
// 			simplifiedProviders = simplifyProvidersRecursively(
// 				simplifiedProviders,
// 				newFullPath,
// 				unknown.Post,
// 			)
// 		}

// 		if unknown.Put.IsNonZero() {
// 			newFullPath := fmt.Sprintf("%s;%s", fullPath, "put")
// 			simplifiedProviders = simplifyProvidersRecursively(
// 				simplifiedProviders,
// 				newFullPath,
// 				unknown.Put,
// 			)
// 		}

// 		if unknown.Delete.IsNonZero() {
// 			newFullPath := fmt.Sprintf("%s;%s", fullPath, "delete")
// 			simplifiedProviders = simplifyProvidersRecursively(
// 				simplifiedProviders,
// 				newFullPath,
// 				unknown.Delete,
// 			)
// 		}

// 		return simplifiedProviders

// 	case GetMethod:
// 		newFullPath := fmt.Sprintf("%s;%s", fullPath, "response")
// 		return simplifyProvidersRecursively(
// 			simplifiedProviders,
// 			newFullPath,
// 			unknown.Responses,
// 		)

// 	case PostMethod:
// 		newFullPath := fmt.Sprintf("%s;%s", fullPath, "response")
// 		simplifiedProviders = simplifyProvidersRecursively(
// 			simplifiedProviders,
// 			newFullPath,
// 			unknown.Responses,
// 		)
// 		if unknown.HasRequestBody() {
// 			newFullPath := fmt.Sprintf("%s;%s", fullPath, "requestBody")
// 			simplifiedProviders = append(
// 				simplifiedProviders,
// 				newSimplifiedProvider(newFullPath, unknown.RequestBody),
// 			)
// 		}
// 		return simplifiedProviders

// 	case PutMethod:
// 		simplifiedProviders = simplifyProvidersRecursively(simplifiedProviders, fmt.Sprintf("%s;%s", fullPath, "response"), unknown.Responses)
// 		if unknown.HasRequestBody() {
// 			newFullPath := fmt.Sprintf("%s;%s", fullPath, "requestBody")
// 			simplifiedProviders = append(
// 				simplifiedProviders,
// 				newSimplifiedProvider(newFullPath, unknown.RequestBody),
// 			)
// 		}
// 		return simplifiedProviders

// 	case DeleteMethod:
// 		newFullPath := fmt.Sprintf("%s;%s", fullPath, "response")
// 		return simplifyProvidersRecursively(
// 			simplifiedProviders,
// 			newFullPath,
// 			unknown.Responses,
// 		)

// 	case map[int]string:
// 		for code, schemaName := range unknown {
// 			newFullPath := fmt.Sprintf("%s;%d", fullPath, code)
// 			simplifiedProviders = append(simplifiedProviders, newSimplifiedProvider(newFullPath, schemaName))
// 		}
// 		return simplifiedProviders

// 	}

// 	return simplifiedProviders
// }
