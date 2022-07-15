/*
Copyright © 2022 Martin Windolph <martin@yoktobit.de>

*/
package cmd

import (
	"errors"
	"log"
	"os"

	hofmod "github.com/hofstadter-io/hof/lib/mod"
	"github.com/spf13/cobra"
	"github.com/yoktobit/yoktogo-flavour/pkg/mod"
)

const (
	LANG               = "cue"
	HOFSTADTER_CUE_DEP = `
require "github.com/hofstadter-io/hof" v0.6.2` + "\n"
	HOFSTADTER_GO_DEP = `
require (
	cuelang.org/go v0.4.3
	github.com/hofstadter-io/hof v0.6.2
	github.com/kirsle/configdir v0.0.0-20170128060238-e45d2f54772f
)` + "\n"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adds a known flavour with custom logic",
	Long: `This command adds a custom flavour with custom logic.
Known flavours:
- cue`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalln("Expected exactly one argument")
			os.Exit(1)
		}
		flavour := args[0]
		switch flavour {
		case LANG:
			moduleName := mod.GetModuleName()
			hofmod.InitLangs()
			err := hofmod.Init(LANG, moduleName)
			if err != nil {
				log.Fatalln(err.Error())
				os.Exit(1)
			}
			err = appendHof()
			if err != nil {
				log.Fatalln(err.Error())
				os.Exit(1)
			}
			err = appendGo()
			if err != nil {
				log.Fatalln(err.Error())
				os.Exit(1)
			}
			err = hofmod.Vendor(LANG)
			if err != nil {
				log.Fatalln(err.Error())
				os.Exit(1)
			}
		}
	},
}

func appendHof() error {
	return appendToFile("cue.mods", HOFSTADTER_CUE_DEP)
}

func appendGo() error {
	return appendToFile("go.mod", HOFSTADTER_GO_DEP)
}

func appendToFile(filename string, text string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return errors.New("could not open " + filename)
	}
	defer file.Close()

	_, err = file.WriteString(text)

	if err != nil {
		return errors.New("could not write text to " + filename)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
