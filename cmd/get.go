/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"

	"github.com/go-git/go-git/v5"
	hofmod "github.com/hofstadter-io/hof/lib/mod"
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
			os.Exit(1)
		}
		err = process(repoDir)
		if err != nil {
			log.Fatalln(err.Error())
			os.Exit(1)
		}
		err = copy.Copy(repoDir, ".", copy.Options{})
		if err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func process(repoDir string) error {

	err := removeDefaultStuff(repoDir)
	if err != nil {
		return err
	}
	err = processDefinition(repoDir)
	if err != nil {
		return err
	}
	err = removeDefinition(repoDir)
	return err
}

func processDefinition(repoDir string) error {

	//schemaBytes := schema.Flavour
	//definitionBytes, err := readFile(repoDir, "yoktogo-flavour.cue")
	//if err != nil {
	//	return err
	//}
	cwd, _ := os.Getwd()
	os.Chdir(repoDir)
	hofmod.InitLangs()
	err := hofmod.Vendor("cue")
	if err != nil {
		return err
	}
	ctx := cuecontext.New()
	filenames := []string{ /*path.Join(repoDir, "schema/Flavour.cue"), */ path.Join(repoDir, "yoktogo-flavour.cue")}
	instances := load.Instances(filenames, nil)
	var root cue.Value
	for _, bi := range instances {
		if bi.Err != nil {
			return bi.Err
		}
		value := ctx.BuildInstance(bi)
		if value.Err() != nil {
			return value.Err()
		}
		err := value.Validate()
		if err != nil {
			return err
		}
		root = value
	}
	//schemaValue := ctx.CompileBytes(schemaBytes)
	//root := ctx.CompileBytes(definitionBytes, cue.Scope(schemaValue))
	//root := ctx.CompileBytes(definitionBytes)
	printAll(root)
	name := root.LookupPath(cue.ParsePath("name"))
	if name.Err() != nil {
		log.Fatalln("cannot lookup name:", name.Err().Error())
		os.Exit(1)
	}
	log.Println("Name of definition", name)
	err = handleExcludes(repoDir, root)
	if err != nil {
		return err
	}

	os.Chdir(cwd)
	return err
}

func handleExcludes(repoDir string, root cue.Value) error {
	excludes := root.LookupPath(cue.ParsePath("excludes"))
	if excludes.Exists() {
		iterator, err := excludes.List()
		if err != nil {
			return err
		}
		for iterator.Next() {
			val := iterator.Value()
			excludePath, err := val.String()
			if err != nil {
				return err
			}
			log.Println("excluding", excludePath)
			if _, err := os.Stat(excludePath); err == nil {
				err = os.RemoveAll(excludePath)
				if err != nil {
					return err
				}
			}
			os.RemoveAll(path.Join(repoDir, excludePath))
		}
	}
	return nil
}

func printAll(v cue.Value) {
	syn := v.Syntax(
		cue.Final(),         // close structs and lists
		cue.Concrete(false), // allow incomplete values
		cue.Definitions(false),
		cue.Hidden(true),
		cue.Optional(true),
		cue.Attributes(true),
		cue.Docs(true),
	)

	// Pretty print the AST, returns ([]byte, error)
	bs, _ := format.Node(
		syn,
		// format.TabIndent(false),
		// format.UseSpaces(2),
	)

	// print to stdout
	fmt.Println(string(bs))
}

func readFile(repoDir string, file string) ([]byte, error) {
	filename := path.Join(repoDir, file)
	info, err := os.Stat(filename)
	if info.IsDir() || err != nil {
		return []byte{}, err
	}
	bytes, err := ioutil.ReadFile(filename)
	return bytes, err
}

func removeDefaultStuff(repoDir string) error {
	err := os.RemoveAll(path.Join(repoDir, ".gitignore"))
	if err != nil {
		return err
	}
	return os.RemoveAll(path.Join(repoDir, ".git"))
}

func removeDefinition(repoDir string) error {
	flavourPath := path.Join(repoDir, "yoktogo-flavour.cue")
	log.Println("excluding default file", flavourPath)
	if _, err := os.Stat(flavourPath); err == nil {
		err = os.RemoveAll(flavourPath)
		if err != nil {
			return err
		}
	}
	return nil
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
