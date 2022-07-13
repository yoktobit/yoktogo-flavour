/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "gets the contents of some repository to add to the current directory",
	Long: `This command clones the linked repository, 
	adds it's contents to the CWD and excludes the git files. 
	After that it does some templating stuff.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			return
		}
		repoLink := args[0]
		repoDir, err := ioutil.TempDir("", "ygf-*")
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Println("Cloning into temp dir", repoDir)
		_, err = git.PlainClone(repoDir, false, &git.CloneOptions{SingleBranch: true, URL: repoLink})
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = copy.Copy(repoDir, ".", copy.Options{})
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
