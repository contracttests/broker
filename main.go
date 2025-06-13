package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/contracttests/broker/internal/dsl"
)

func main() {
	content := []byte(`
    {
      "api": {
        "name": "catalog-service"
      },
      "provides": {
        "rest": {
          "/products": {
            "get": {
              "responses": {
                "200": "Products"
              }
            },
            "post": {
              "requestBody": "CreateProduct",
              "responses": {
                "201": "Product",
                "400": "BadRequest"
              }
            }
          },
          "/products/{uuid}": {
            "get": {
              "responses": {
                "200": "Product",
                "404": "NotFound"
              }
            },
            "put": {
              "requestBody": "UpdateProduct",
              "responses": {
                "200": "Product",
                "400": "BadRequest",
                "404": "NotFound"
              }
            }
          }
        }
      },
      "consumes": {
        "payments-service": {
          "rest": {
            "/payments": {
              "post": {
                "requestBody": "CreatePayment",
                "responses": {
                  "201": "Payment",
                  "400": "BadRequest"
                }
              }
            }
          }
        },
        "users-service": {
          "rest": {
            "/users/{uuid}": {
              "get": {
                "responses": {
                  "200": "User",
                  "404": "NotFound"
                }
              }
            }
          }
        }
      },
      "schemas": {
        "Product": {
          "type": "object",
          "properties": {
            "uuid": {
              "type": "string"
            },
            "name": {
              "type": "string"
            },
            "description": {
              "type": "string"
            },
            "price": {
              "type": "number"
            },
            "stock": {
              "type": "integer"
            },
            "category": {
              "type": "string"
            }
          }
        },
        "Products": {
          "type": "array",
          "items": {
            "$ref": "Product"
          }
        },
        "CreateProduct": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            },
            "description": {
              "type": "string"
            },
            "price": {
              "type": "number"
            },
            "stock": {
              "type": "integer"
            },
            "category": {
              "type": "string"
            }
          }
        },
        "UpdateProduct": {
          "type": "object",
          "properties": {
            "uuid": {
              "type": "string"
            },
            "name": {
              "type": "string"
            },
            "description": {
              "type": "string"
            },
            "price": {
              "type": "number"
            },
            "stock": {
              "type": "integer"
            },
            "category": {
              "type": "string"
            }
          }
        },
        "User": {
          "type": "object",
          "properties": {
            "uuid": {
              "type": "string"
            },
            "name": {
              "type": "string"
            },
            "email": {
              "type": "string"
            },
            "address": {
              "type": "string"
            }
          }
        },
        "Payment": {
          "type": "object",
          "properties": {
            "uuid": {
              "type": "string"
            },
            "userId": {
              "type": "string"
            },
            "orderId": {
              "type": "string"
            },
            "amount": {
              "type": "number"
            },
            "status": {
              "type": "string"
            }
          }
        },
        "CreatePayment": {
          "type": "object",
          "properties": {
            "userId": {
              "type": "string"
            },
            "orderId": {
              "type": "string"
            },
            "amount": {
              "type": "number"
            }
          }
        },
        "BadRequest": {
          "type": "object",
          "properties": {
            "message": {
              "type": "string"
            }
          }
        },
        "NotFound": {
          "type": "object",
          "properties": {
            "message": {
              "type": "string"
            }
          }
        }
      }
    }
	`)

	contract := dsl.Contract{}

	if err := json.Unmarshal(content, &contract); err != nil {
		log.Fatal("could not parse json, make sure it is valid contract file")
	}

	comparableContract := dsl.NewComparableContract(contract)

	for _, comparableResource := range comparableContract.ConsumesResources {
		fmt.Printf("%s -> %s\n", comparableResource.Path, comparableResource.SchemaName)
	}

	for _, comparableResource := range comparableContract.ProvidesResources {
		fmt.Printf("%s -> %s\n", comparableResource.Path, comparableResource.SchemaName)
	}

	for _, schemaName := range comparableContract.Schemas {
		for _, property := range schemaName.Properties {
			fmt.Printf("%s -> %s -> %s\n", schemaName.Name, property.Path, property.Type)
		}
	}
}
