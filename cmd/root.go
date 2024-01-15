package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thecheerfuldev/gitcd-go/config"
	"github.com/thecheerfuldev/gitcd-go/repository"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitcd [git repo]",
	Args:  cobra.MaximumNArgs(1),
	Short: "",
	Long: `GitCD is a CLI tool that lets you easily index and navigate to git projects.
If you don't provide a repo to search for, a top 10 will be displayed.'`,
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
			repository.GiveTopTen()
			handleMultipleMatches(repository.GiveTopTen())
			return
		}

		if len(args) == 1 {
			matches := repository.GetProjectContaining(args[0])
			if len(matches) == 1 {
				handleSingleMatch(matches[0])
				return
			}
			handleMultipleMatches(matches)
		}
	},
}

func handleSingleMatch(match string) {
	err := os.WriteFile(config.GetDirChangerPath(), generateCdScript(match), 0755)
	if err != nil {
		fmt.Println("Something went wrong while preparing to change directory:", err)
		_ = os.Remove(config.GetDirChangerPath())
		os.Exit(1)
	}
	project := repository.GetProject(match)
	project.UpdateCounter()
	fmt.Println("Changing directory to:", match)
}

func handleMultipleMatches(matches []string) {
	var numberOfOptions int64 = 1
	for range matches {
		fmt.Printf("%v) %v\n", numberOfOptions, matches[numberOfOptions-1])
		numberOfOptions++
	}

	fmt.Print("Pick a project: ")
	var choice string
	_, _ = fmt.Scan(&choice)

	convertedChoice, err := strconv.ParseInt(choice, 10, 0)

	if err != nil || convertedChoice > numberOfOptions {
		fmt.Println("Invalid choice.")
		return
	}

	if convertedChoice == 0 {
		return
	}

	err = os.WriteFile(config.GetDirChangerPath(), generateCdScript(matches[convertedChoice-1]), 0755)
	if err != nil {
		fmt.Println("Something went wrong while preparing to change directory:", err)
		_ = os.Remove(config.GetDirChangerPath())
		os.Exit(1)
	}
	project := repository.GetProject(matches[convertedChoice-1])
	project.UpdateCounter()
	fmt.Println("Changing directory to:", matches[convertedChoice-1])
}

func generateCdScript(path string) []byte {
	return []byte(fmt.Sprintf(
		`#! /bin/bash
cd %v`, path))
}

func handleScanFlag() {
	if _, err := os.Stat(config.GetProjectRootPath()); os.IsNotExist(err) {
		fmt.Println("$GITCD_PROJECT_HOME does not exist")
	}

	filepath.WalkDir(config.GetProjectRootPath(), func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == ".git" {
			parentDir := strings.Replace(path, "/.git", "", 1)
			repository.AddProject(parentDir, 0)
		}
		return nil
	})
}

func handleCleanFlag() {
	for _, path := range repository.GetAllProjects() {
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
