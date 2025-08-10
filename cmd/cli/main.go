package main

import (
	"fmt"
	"os"

	"github.com/contracttests/broker/internal/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(cmd.ValidateContractCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
