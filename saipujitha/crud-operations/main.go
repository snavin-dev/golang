package main

import (
	"log"

	"github.com/saipujitha/crud-operations/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
