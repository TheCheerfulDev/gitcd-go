# GitCD

GitCD is a CLI tool that lets you easily index and navigate to git projects.

# Prerequisites

* MacOS or Linux
* [homebrew](https://brew.sh/) installed

# Installation

```bash
brew tap thecheerfuldev/cli
brew install thecheefuldev/cli/gitcd
```

Then add the following to your profile (.profile .bashrc .zshrc etc...)

```bash
export GITCD_PROJECT_HOME=</your/projects/root>
alias alias gcd="source gitcd-go-runner.sh"
```

Or pick any alias that you prefer.

NOTE: If GITCD_PROJECT_HOME is NOT set, your home directory will be used instead.

## Automate Scanning

If you wish to automate scanning for new projects, you can add the following cron job (or similar):

```text
0 7,12,15,20 * * * gcd --scan > /dev/null # scan for projects at 7am, 12pm, 3pm and 8pm
```

# Usage

You can find all commands with

```bash
gcd --help
```

## Scanning

Before your first usage, you should can for git repositories

```bash
gcd --scan
```

## Cleaning Database

Purge repositories that no longer exist

```bash
gcd --clean
```

## Full Reset

Clear the entire database and scan for git repositories

```bash
gcd --reset
```

## Environment variables

* GITCD_PROJECT_HOME - Root directory for your projects

# Special thanks

* [IvoNet](https://github.com/IvoNet) for creating the original version of this tool, and pushing me to rewrite it
  in Go!