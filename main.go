/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/thecheerfuldev/gitcd/cmd"
	"github.com/thecheerfuldev/gitcd/repository"
)

func main() {
	repository.Init()
	defer repository.WriteChangesToDatabase()
	cmd.Execute()
}
