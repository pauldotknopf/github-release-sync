# What?

A helper tool (written in Go) that helps synchronize a GitHub release's files with a local directory.

This is useful for storing Arch Linux or Void Linux repositories on GitHub (see [here](https://github.com/pauldotknopf/binary-repos/releases/tag/void-linux-release) and [here](https://github.com/pauldotknopf/binary-repos/releases/tag/arch-linux-release)).

# Installation

```
go get github.com/pauldotknopf/github-release-sync
```

# Usage

```
NAME:
   github-release-sync - Sync a directory of files with a GitHub release.

USAGE:
   github-release-sync [command options]

VERSION:
   0.0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --access-token value  the github access token to use the api, $GITHUB_ACCESS_TOKEN if empty
   --owner value         the owner of the repository (user/org)
   --repo value          the repository
   --tag value           the tagged release
   --overwrite           overwrite files that already exist on github
   --prune               delete release assets on the remote that don't exist locally
   --dir value           the directory to sync to the github release (default: ".")
   --help, -h            show help
   --version, -v         print the version
```
