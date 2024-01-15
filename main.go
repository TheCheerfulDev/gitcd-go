/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"github.com/thecheerfuldev/gitcd-go/cmd"
	"github.com/thecheerfuldev/gitcd-go/config"
	"github.com/thecheerfuldev/gitcd-go/repository"
	"time"
)

func main() {
	start := time.Now()
	config.Init()
	repository.Init()
	defer repository.WriteChangesToDatabase()
	cmd.Execute()
	defer fmt.Printf("Took %v\n", time.Since(start))

}
