package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thecheerfuldev/gitcd-go/config"
	"github.com/thecheerfuldev/gitcd-go/repository"
	"github.com/theckman/yacspin"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "gitcd [git repo]",
	//Args:    cobra.MaximumNArgs(1),
	Version: "1.0.4",
	Short:   "",
	Long: `GitCD is a CLI tool that lets you easily index and navigate to git projects.
If you don't provide a repo to search for, a top 10 will be displayed.'`,
	Run: func(cmd *cobra.Command, args []string) {
		resetFlagUsed, _ := cmd.Flags().GetBool("reset")
		if resetFlagUsed {
			repository.ResetDatabase()
			handleScanFlag()
			return
		}

		scanFlagUsed, _ := cmd.Flags().GetBool("scan")
		if scanFlagUsed {
			handleScanFlag()
			return
		}

		if len(repository.GetAllProjects()) == 0 {
			fmt.Println("Your database appears to be empty. Run gitcd with the --scan flag to index your git projects.")
			os.Exit(1)
		}

		cleanFlagUsed, _ := cmd.Flags().GetBool("clean")
		if cleanFlagUsed {
			handleCleanFlag()
			return
		}

		if len(args) == 0 {
			handleMultipleMatches(repository.GiveTopTen())
			return
		}

		if len(args) > 0 {
			expression := extractExpression(args)
			matches := repository.GetProjectsRegex(expression)
			if len(matches) == 0 {
				fmt.Println("No projects found.")
				os.Exit(0)
			}

			if len(matches) == 1 {
				handleSingleMatch(matches[0])
				return
			}
			handleMultipleMatches(matches)
		}
	},
}

func extractExpression(args []string) string {
	var expression = args[0]

	for i := 1; i < len(args); i++ {
		expression += ".*" + args[i]
	}
	return expression

}

func handleSingleMatch(match string) {
	err := os.WriteFile(config.Get().DirChangerPath, generateCdScript(match), 0755)
	if err != nil {
		fmt.Println("Something went wrong while preparing to change directory:", err)
		_ = os.Remove(config.Get().DirChangerPath)
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

	fmt.Print("Select a project: ")
	var choice string
	_, _ = fmt.Scan(&choice)

	if choice == "q" || choice == "0" {
		fmt.Println("Quitting.")
		return
	}

	convertedChoice, err := strconv.ParseInt(choice, 10, 0)

	if err != nil || convertedChoice > numberOfOptions {
		fmt.Println("Invalid choice.")
		return
	}

	handleSingleMatch(matches[convertedChoice-1])
}

func generateCdScript(path string) []byte {
	shell, b := os.LookupEnv("SHELL")
	if !b {
		shell = "/bin/bash"
	}
	return []byte(fmt.Sprintf(
		`#!%s
cd %s`, shell, path))
}

func handleScanFlag() {
	root := config.Get().ProjectRootPath
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Printf("$GITCD_PROJECT_HOME (%s) does not exist\n", root)
		return
	}

	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[2],
		Colors:          []string{"fgYellow"},
		Suffix:          " Scanning for git projects, this might take a while...",
		SuffixAutoColon: false,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
		StopMessage:     " Done!",
	}
	s, _ := yacspin.New(cfg)
	_ = s.Start()

	fmt.Printf("Scanning %v for git projects. This might take a while... ", root)
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == ".git" {
			repository.AddProject(filepath.Dir(path))
		}
		return nil
	})
	_ = s.Stop()
}

func handleCleanFlag() {
	for _, path := range repository.GetAllProjects() {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			removeProject(path)
			continue
		}

		if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
			removeProject(path)
		}
	}
}

func removeProject(path string) {
	repository.RemoveProject(path)
	fmt.Println("Removed:", path)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("gitcd version %s - © Mark Hendriks <thecheerfuldev>\n", rootCmd.Version))
	rootCmd.Flags().BoolP("scan", "", false, "Scan for git projects in $GITCD_PROJECT_HOME")
	rootCmd.Flags().BoolP("clean", "", false, "Remove all projects that no longer exist")
	rootCmd.Flags().BoolP("reset", "", false, "Resets the database and scans for git project in $GITCD_PROJECT_HOME")
}
