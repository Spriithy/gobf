package main

import (
	"log"
	"os"

	"github.com/Spriithy/gobf"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		vm := gobf.NewInterpreterFromFile(args[1])
		vm.Exec()
		return
	}

	log.Fatalf("no input file")
}
