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
	Use:   "gitcd [git repo]",
	Args:  cobra.MaximumNArgs(1),
	Short: "",
	Long:  `GitCD is a CLI tool that lets you easily index and navigate to git projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		scanFlagUsed, _ := cmd.Flags().GetBool("scan")
		if scanFlagUsed {
			handleScanFlag()
			return
		}
		cleanFlagUsed, _ := cmd.Flags().GetBool("clean")
		if cleanFlagUsed {
			handleCleanFlag()
			return
		}

		if len(args) == 0 {

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
			fmt.Scan(&choice)

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

func handleScanFlag() {
	filepath.WalkDir("/Users/mark/dev", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == ".git" {
			parentDir := strings.Replace(path, "/.git", "", 1)
			repository.AddProject(parentDir, 0)
		}
		return nil
	})
}

func handleCleanFlag() {
	for _, path := range repository.GetAllProjects() {
		// TODO check if path still exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			repository.RemoveProject(path)
			fmt.Println("Removed:", path)
		}
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("scan", "", false, "Scan for git projects in $GITCD_PROJECT_HOME")
	rootCmd.Flags().BoolP("clean", "", false, "Remove all projects that no longer exist from GITCD")
}
