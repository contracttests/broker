package tool

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func UploadContract(contractFilepath string) {
	destinationFilepath := "./contracts/dsl-validated/"
	contractFileName := strings.TrimSuffix(filepath.Base(contractFilepath), ".yaml")
	contractContent, err := os.ReadFile(contractFilepath)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(fmt.Sprintf("%s%s.yaml", destinationFilepath, contractFileName), contractContent, 0644)
}
