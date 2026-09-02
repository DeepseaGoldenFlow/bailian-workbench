package main

import (
	"encoding/json"
	"log"
	"os"

	"bailian-workbench/internal/catalog"
)

func main() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(catalog.BuildCatalog()); err != nil {
		log.Fatal(err)
	}
}
