package cmd

import (
	"fmt"
	"github.com/thecheerfuldev/gitcd/repository"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitcd",
	Args:  cobra.MaximumNArgs(1),
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		scanFlagUsed, _ := cmd.Flags().GetBool("scan")
		if scanFlagUsed {
			walkDirectoryTree()
			return
		}

		if len(args) == 1 {

			matches := repository.GetProjectContaining(args[0])

			if len(matches) == 1 {
				os.WriteFile("/tmp/gitcd-action.sh", generateCdScript(matches[0]), 0755)
				project := repository.GetProject(matches[0])
				project.UpdateCounter()
				fmt.Println("Changing directory to:", matches[0])
				return
			}
			// found multiple, so make a table with numbers
			fmt.Println("Found multiple:", matches)

			var numberOfOptions int64 = 1
			for range matches {
				fmt.Printf("%v) %v\n", numberOfOptions, matches[numberOfOptions-1])
				numberOfOptions++
			}

			fmt.Print("Pick a project: ")
			var choice string
			fmt.Scanln(&choice)
			choice = strings.ReplaceAll(choice, "\n", "")
			choice = strings.ReplaceAll(choice, "\r", "")

			convertedChoice, err := strconv.ParseInt(choice, 10, 0)

			if err != nil || convertedChoice > numberOfOptions {
				fmt.Println("Invalid choice.")
				return
			}

			if convertedChoice == 0 {
				return
			}

			os.WriteFile("/tmp/gitcd-action.sh", generateCdScript(matches[convertedChoice-1]), 0755)
			project := repository.GetProject(matches[convertedChoice-1])
			project.UpdateCounter()
			fmt.Println("Changing directory to:", matches[convertedChoice-1])
		}
	},
}

func generateCdScript(path string) []byte {
	return []byte(fmt.Sprintf(
		`#! /bin/bash
cd %v`, path))
}

func walkDirectoryTree() {
	filepath.WalkDir("/Users/mark/dev", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == ".git" {
			parentDir := strings.Replace(path, "/.git", "", 1)
			//fmt.Println("Found .git directory in: ", parentDir)
			split := strings.Split(parentDir, "/")
			//fmt.Println(split[len(split)-1])
			repository.AddProject(split[len(split)-1], parentDir, 0)
		}
		return nil
	})
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitcd.yaml)")

	// Cobra also supports local flags, which will only gitcd-go-runner.sh
	// when this action is called directly.
	rootCmd.Flags().BoolP("scan", "s", false, "Scan for git projects in $GITCD_PROJECT_HOME")
}
