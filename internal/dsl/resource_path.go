package dsl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/contracttests/broker/server/internal/model"
)

var (
    consumerRestRequestRegex = regexp.MustCompile(
        `^(?P<owner>[^;]+);consumes;(?P<provider>[^;]+);rest;(?P<endpoint>[^;]+);(?P<method>[^;]+);request$`,
    )
    consumerRestResponseRegex = regexp.MustCompile(
        `^(?P<owner>[^;]+);consumes;(?P<provider>[^;]+);rest;(?P<endpoint>[^;]+);(?P<method>[^;]+);responses;(?P<status>\d+)$`,
    )
    providerRestRequestRegex = regexp.MustCompile(
        `^(?P<owner>[^;]+);provides;rest;(?P<endpoint>[^;]+);(?P<method>[^;]+);request$`,
    )
    providerRestResponseRegex = regexp.MustCompile(
        `^(?P<owner>[^;]+);provides;rest;(?P<endpoint>[^;]+);(?P<method>[^;]+);responses;(?P<status>\d+)$`,
    )
)

type ResourcePath string

func NewResourcePath(resourcePath string) ResourcePath {
	return ResourcePath(resourcePath)
}

func (rp *ResourcePath) Append(parts ...string) ResourcePath {
	separator := ";"

	if string(*rp) == "" {
		return ResourcePath(strings.Join(parts, separator))
	}

	chunks := strings.Join(parts, separator)

	return ResourcePath(strings.Join([]string{string(*rp), chunks}, separator))
}

func (rp *ResourcePath) String() string {
	return string(*rp)
}

func (rp *ResourcePath) Split() []string {
	return strings.Split(rp.String(), ";")
}

func (rp *ResourcePath) IsConsumer() bool {
	return strings.Contains(rp.String(), "consumes")
}

func (rp *ResourcePath) IsProvider() bool {
	return strings.Contains(rp.String(), "provides")
}

func (rp *ResourcePath) ExtractNamedArgs(regex *regexp.Regexp) (map[string]string, bool) {
    match := regex.FindStringSubmatch(rp.String())
    if match == nil {
        return nil, false
    }

    args := make(map[string]string, len(regex.SubexpNames()))
    for i, name := range regex.SubexpNames() {
        if name == "" {
			continue
        }

        args[name] = match[i]
    }

    return args, true
}

func (rp *ResourcePath) ToConsumerRestRequestArgs() model.ConsumerRestRequestArgs {
	args, ok := rp.ExtractNamedArgs(consumerRestRequestRegex)
	if !ok {
		panic(fmt.Errorf("invalid consumer rest request path: %q", rp.String()))
	}

	return model.ConsumerRestRequestArgs{
		Owner:    args["owner"],
		Provider: args["provider"],
		Endpoint: args["endpoint"],
		Method:   args["method"],
	}
}

func (rp *ResourcePath) ToConsumerRestResponseArgs() model.ConsumerRestResponseArgs {
    args, ok := rp.ExtractNamedArgs(consumerRestResponseRegex)
    if !ok {
        panic(fmt.Errorf("invalid consumer rest response path: %q", rp.String()))
    }

    return model.ConsumerRestResponseArgs{
        Owner:      args["owner"],
        Provider:   args["provider"],
        Endpoint:   args["endpoint"],
        Method:     args["method"],
        StatusCode: args["status"],
    }
}

func (rp *ResourcePath) ToProviderRestRequestArgs() model.ProviderRestRequestArgs {
	args, ok := rp.ExtractNamedArgs(providerRestRequestRegex)
	if !ok {
		panic(fmt.Errorf("invalid provider rest request path: %q", rp.String()))
	}

	return model.ProviderRestRequestArgs{
		Owner:    args["owner"],
		Endpoint: args["endpoint"],
		Method:   args["method"],
	}
}

func (rp *ResourcePath) ToProviderRestResponseArgs() model.ProviderRestResponseArgs {
	args, ok := rp.ExtractNamedArgs(providerRestResponseRegex)
	if !ok {
		panic(fmt.Errorf("invalid provider rest response path: %q", rp.String()))
	}

	return model.ProviderRestResponseArgs{
		Owner:      args["owner"],
		Endpoint:   args["endpoint"],
		Method:     args["method"],
		StatusCode: args["status"],
	}
}
