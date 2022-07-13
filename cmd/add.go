/*
Copyright © 2022 Martin Windolph <martin@yoktobit.de>

*/
package cmd

import (
	"log"

	hofmod "github.com/hofstadter-io/hof/lib/mod"
	"github.com/spf13/cobra"
	"github.com/yoktobit/yoktogo-flavour/pkg/mod"
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
			return
		}
		flavour := args[0]
		switch flavour {
		case "cue":
			moduleName := mod.GetModuleName()
			hofmod.InitLangs()
			err := hofmod.Init("cue", moduleName)
			if err != nil {
				log.Fatalln(err.Error())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
