package mod

import (
	"io/ioutil"
	"log"

	modfile "golang.org/x/mod/modfile"
)

func GetModuleName() (string, error) {
	goModBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		log.Fatalln(err.Error())
		return "", err
	}
	return modfile.ModulePath(goModBytes), nil
}
