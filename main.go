package main

import (
	"fmt"
	"github.com/thecheerfuldev/gitcd-go/cmd"
	"github.com/thecheerfuldev/gitcd-go/config"
	"github.com/thecheerfuldev/gitcd-go/repository"
	"os"
)

func main() {
	c := config.Default()
	err := config.Init(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	repository.Init(c)
	defer repository.WriteChangesToDatabase()
	cmd.Execute()
}
