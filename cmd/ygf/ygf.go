package main

import (
	"log"

	"github.com/yoktobit/yoktogo-flavour/pkg/mod"
)

func main() {
	modName := mod.GetModuleName()
	log.Println("modName=", modName)
}
