package main

import (
	"github.com/thecheerfuldev/gitcd-go/cmd"
	"github.com/thecheerfuldev/gitcd-go/config"
	"github.com/thecheerfuldev/gitcd-go/repository"
)

func main() {
	config.Init()
	repository.Init()
	defer repository.WriteChangesToDatabase()
	cmd.Execute()
}
