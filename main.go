package main

import (
	"fmt"
	"os"

	"github.com/thecheerfuldev/gitcd-go/cmd"
	"github.com/thecheerfuldev/gitcd-go/config"
	"github.com/thecheerfuldev/gitcd-go/repository"
)

func main() {
	c := config.Default()
	err := config.Init(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = repository.Init(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer repository.WriteChangesToDatabase()
	cmd.Execute()
}
