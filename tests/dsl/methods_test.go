package dsl_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/contracttests/broker/internal/dsl"
	"github.com/stretchr/testify/assert"
)

func TestGetMethod(t *testing.T) {
	content := []byte(`
		{
			"responses": {
				"200": "Products"
			}
		}
	`)

	actual := dsl.GetMethod{}

	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.GetMethod{
		Responses: dsl.Responses{
			200: "Products",
		},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
	assert.True(t, actual.IsNonZero())
}

func TestGetMethodWithEmptyResponses(t *testing.T) {
	content := []byte(`
		{
			"responses": {}
		}
	`)

	actual := dsl.GetMethod{}

	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.GetMethod{
		Responses: dsl.Responses{},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
	assert.False(t, actual.IsNonZero())
}

func TestPostMethod(t *testing.T) {
	content := []byte(`
		{
			"requestBody": "CreateProduct",
			"responses": {
				"201": "Product",
				"400": "BadRequest"
			}
		}
	`)

	actual := dsl.PostMethod{}

	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.PostMethod{
		RequestBody: "CreateProduct",
		Responses: dsl.Responses{
			201: "Product",
			400: "BadRequest",
		},
	}

	assert.ObjectsAreEqual(expected, actual)
	assert.True(t, actual.HasRequestBody())
	assert.True(t, actual.IsNonZero())
}

func TestZeroPostMethod(t *testing.T) {
	content := []byte(`
		{
			"responses": {}
		}
	`)

	actual := dsl.PostMethod{}
	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.PostMethod{
		RequestBody: "",
		Responses:   dsl.Responses{},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
	assert.Equal(t, 0, len(actual.Responses))
}

func TestPutMethod(t *testing.T) {
	content := []byte(`
		{
			"requestBody": "UpdateProduct",
			"responses": {
				"200": "Product",
				"400": "BadRequest",
				"404": "NotFound"
			}
		}
	`)

	actual := dsl.PutMethod{}

	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.PutMethod{
		RequestBody: "UpdateProduct",
		Responses: dsl.Responses{
			200: "Product",
			400: "BadRequest",
			404: "NotFound",
		},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
	assert.True(t, actual.IsNonZero())
}

func TestZeroPutMethod(t *testing.T) {
	content := []byte(`
		{
			"responses": {}
		}
	`)

	actual := dsl.PutMethod{}
	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.PutMethod{
		RequestBody: "",
		Responses:   dsl.Responses{},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
	assert.False(t, actual.IsNonZero())
}

func TestDeleteMethod(t *testing.T) {
	content := []byte(`
		{
			"responses": {
				"201": "Unknown",
				"404": "NotFound"
			}
		}
	`)

	actual := dsl.DeleteMethod{}

	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.DeleteMethod{
		Responses: dsl.Responses{
			201: "Unknown",
			404: "NotFound",
		},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
	assert.True(t, actual.IsNonZero())
}

func TestZeroDeleteMethod(t *testing.T) {
	content := []byte(`
		{
			"responses": {}
		}
	`)

	actual := dsl.DeleteMethod{}
	if err := json.Unmarshal(content, &actual); err != nil {
		log.Fatal(err)
	}

	expected := dsl.DeleteMethod{
		Responses: dsl.Responses{},
	}

	assert.True(t, assert.ObjectsAreEqual(expected, actual))
	assert.False(t, actual.IsNonZero())
}
