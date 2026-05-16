package upload_contract_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/contracttests/broker/server/internal"
	"github.com/contracttests/broker/server/internal/components"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	Components *components.Components
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) SetupTest() {
	suite.Components = internal.Run()
}

type Request struct {
	Method string
	Path string
	Body string
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body string
}

func (suite *Suite) Request(args Request) (*Response, error) {
	request := httptest.NewRequest(args.Method, args.Path, strings.NewReader(args.Body))

	for key, value := range args.Headers {
		request.Header.Set(key, value)
	}

	response, err := suite.Components.Server.Test(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: response.StatusCode,
		Body: string(body),
	}, nil
}
