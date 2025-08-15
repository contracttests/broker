package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/contracttests/broker/internal/dsl"
	"github.com/contracttests/broker/internal/tool"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var ValidateContractCmd = &cobra.Command{
	Use:  "contract:validate",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tool.LoadValidatedContracts()

		contractDslFilepath := args[0]

		contractDslContent, err := os.ReadFile(contractDslFilepath)
		if err != nil {
			log.Fatal(err)
		}

		var contractDsl dsl.Contract
		if err := yaml.Unmarshal(contractDslContent, &contractDsl); err != nil {
			log.Fatal(err)
		}

		if !isContractFileNameValid(contractDslFilepath, contractDsl) {
			log.Fatalf("Contract filename must match with the contract name: %s", contractDsl.Api.Name)
		}

		contractModel := contractDsl.ToContractModel()

		// content, err := json.Marshal(contractModel)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println(string(content))
		// os.Exit(0)

		hasError := tool.ValidateContract(contractModel)

		if hasError {
			os.Exit(1)
		}

		fmt.Println("Contract is valid")
	},
}

func isContractFileNameValid(contractDslFilepath string, contractDsl dsl.Contract) bool {
	contractDslFileName := strings.TrimSuffix(filepath.Base(contractDslFilepath), ".yaml")
	return contractDsl.Api.Name == contractDslFileName
}
