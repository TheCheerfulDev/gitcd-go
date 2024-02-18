# GitCD

GitCD is a CLI tool that lets you easily index and navigate to git projects.

# Prerequisites

* MacOS or Linux
* Optional: [homebrew](https://brew.sh/) installed
* Optional: [Go 1.22+](https://golang.org/dl/) installed

# Installation

## Homebrew

1. Tap the repository

    ```bash
    brew tap thecheerfuldev/cli
    ```
2. Install gitcd

    ```bash
    brew install thecheerfuldev/cli/gitcd
    ```

3. Add the following to your profile (.profile .bashrc .zshrc etc...)

    ```bash
    export GITCD_PROJECT_HOME=</your/projects/root>
    alias gcd="source gitcd-go-runner.sh"
    ```

Or pick any alias that you prefer.

NOTE: If GITCD_PROJECT_HOME is NOT set, your home directory will be used instead.

## Download and install manually

1. Download the latest release from the [releases page](https://github.com/TheCheerfulDev/gitcd-go/releases/latest)
2. Extract the archive
3. Move the binary and runner script to a location in your PATH
4. Add the following to your profile (.profile .bashrc .zshrc etc...)

    ```bash
    export GITCD_PROJECT_HOME=</your/projects/root>
    alias gcd="source gitcd-go-runner.sh"
    ```

Or pick any alias that you prefer.

NOTE: If GITCD_PROJECT_HOME is NOT set, your home directory will be used instead.

## Build from source

1. Clone the repository to a location of your choice

    ```bash
    git clone https://github.com/TheCheerfulDev/gitcd-go.git
    ```
2. Make sure the location of the cloned repository is in your PATH
3. Build the binary

    ```bash
    go build
    ```

4. Add the following to your profile (.profile .bashrc .zshrc etc...)

    ```bash
    export GITCD_PROJECT_HOME=</your/projects/root>
    alias gcd="source gitcd-go-runner.sh"
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

### Scanning

Before your first usage, you should can for git repositories

```bash
gcd --scan
```

### Searching

Search for a project

```bash
gcd <part of project name>
```

You can also use regex, make sure to use quotes

```bash
gcd ".*project.*"
```

You can also use multiple search terms, separated by spaces. Each term is stitched together with the .* regex

```bash
gcd first second third
```

### Cleaning Database

Purge repositories that no longer exist

```bash
gcd --clean
```

### Full Reset

Clear the entire database and scan for git repositories

```bash
gcd --reset
```

### Environment variables

* GITCD_PROJECT_HOME - Root directory for your projects

# License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details

# Special thanks

* [IvoNet](https://github.com/IvoNet) for creating the original version of this tool, and pushing me to rewrite it
  in Go!