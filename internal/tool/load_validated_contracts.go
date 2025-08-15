package tool

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/contracttests/broker/internal/dsl"
	"github.com/contracttests/broker/internal/model"
	"github.com/ghodss/yaml"
)

type ContractFilepath struct {
	SpecFilepath       string
	ModelCacheFilepath string
}

func getContractFilepaths() []ContractFilepath {
	dslValidatedFilepaths, err := os.ReadDir("contracts/dsl-validated")
	if err != nil {
		log.Fatal(err)
	}

	contractFilepaths := make([]ContractFilepath, 0)
	for _, contractFilepath := range dslValidatedFilepaths {
		if contractFilepath.IsDir() {
			continue
		}

		contractFileName := strings.TrimSuffix(filepath.Base(contractFilepath.Name()), ".yaml")

		contractFilepaths = append(contractFilepaths, ContractFilepath{
			SpecFilepath:       fmt.Sprintf("contracts/dsl-validated/%s.yaml", contractFileName),
			ModelCacheFilepath: fmt.Sprintf("contracts/model-cache/%s.json", contractFileName),
		})
	}

	return contractFilepaths
}

func LoadValidatedContracts() {
	contractFilepaths := getContractFilepaths()

	for _, contractFilepath := range contractFilepaths {
		if _, err := os.Stat(contractFilepath.ModelCacheFilepath); err == nil {
			contractModelCacheContent, err := os.ReadFile(contractFilepath.ModelCacheFilepath)
			if err != nil {
				log.Fatal(err)
			}

			var contractModel model.Contract
			json.Unmarshal(contractModelCacheContent, &contractModel)

			SaveContract(contractModel)
			continue
		}

		var contractDsl dsl.Contract
		contractDslContent, err := os.ReadFile(contractFilepath.SpecFilepath)
		if err != nil {
			log.Fatal(err)
		}

		if err := yaml.Unmarshal(contractDslContent, &contractDsl); err != nil {
			log.Fatal(err)
		}

		contractModel := contractDsl.ToContractModel()
		SaveContract(contractModel)

		contractModelCacheContent, err := json.Marshal(contractModel)
		if err != nil {
			log.Fatal(err)
		}

		os.MkdirAll(filepath.Dir(contractFilepath.ModelCacheFilepath), 0755)
		os.WriteFile(contractFilepath.ModelCacheFilepath, contractModelCacheContent, 0644)
	}
}
