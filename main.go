package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/go-github/github"

	"github.com/urfave/cli"
)

func main() {
	var accessToken string
	var owner string
	var repo string
	var tag string
	var overwrite bool
	var prune bool
	var dir string

	app := cli.NewApp()
	app.Name = "darch"
	app.Usage = "A tool used to build, boot and share stateless Arch images."
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "access-token",
			Usage:       "the github access token to use the api, $GITHUB_ACCESS_TOKEN if empty",
			Destination: &accessToken,
		},
		cli.StringFlag{
			Name:        "owner",
			Usage:       "the owner of the repository (user/org)",
			Destination: &owner,
		},
		cli.StringFlag{
			Name:        "repo",
			Usage:       "the repository",
			Destination: &repo,
		},
		cli.StringFlag{
			Name:        "tag",
			Usage:       "the tagged release",
			Destination: &tag,
		},
		cli.BoolFlag{
			Name:        "overwrite",
			Usage:       "overwrite files that already exist on github",
			Destination: &overwrite,
		},
		cli.BoolFlag{
			Name:        "prune",
			Usage:       "delete release assets on the remote that don't exist locally",
			Destination: &prune,
		},
		cli.StringFlag{
			Name:        "dir",
			Usage:       "the directory to sync to the github release",
			Value:       ".",
			Destination: &dir,
		},
	}
	app.Action = func(clicontext *cli.Context) error {
		if len(accessToken) == 0 {
			accessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
			if len(accessToken) == 0 {
				return fmt.Errorf("access token is required")
			}
		}
		if len(owner) == 0 {
			return fmt.Errorf("owner is required")
		}
		if len(repo) == 0 {
			return fmt.Errorf("repo is required")
		}
		if len(tag) == 0 {
			return fmt.Errorf("tag is required")
		}
		if len(dir) == 0 {
			return fmt.Errorf("directory is required")
		}
		dir = ExpandPath(dir)
		if stat, err := os.Stat(dir); err != nil {
			return err
		} else if !stat.IsDir() {
			return fmt.Errorf("%s is not a directory", dir)
		}

		fmt.Println("getting release...")
		release, err := GetReleasesTag(owner, repo, tag)
		if err != nil {
			return err
		}

		fmt.Println("getting current assets...")
		assets, err := GetReleaseAssets(owner, repo, release.GetID())
		if err != nil {
			return err
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			var releaseAsset *github.ReleaseAsset
			for _, asset := range assets {
				if asset.GetName() == file.Name() {
					releaseAsset = asset
				}
			}
			if releaseAsset != nil {
				fmt.Printf("%s already uploaded\n", file.Name())
				//continue
				if overwrite {
					fmt.Println("overwriting...")
					fmt.Printf("deleting existing asset %s...\n", file.Name())
					err = DeleteAsset(owner, repo, releaseAsset.GetID(), accessToken)
					if err != nil {
						return err
					}
				} else {
					continue
				}
			}
			fmt.Printf("uploading %s\n", file.Name())
			err = UploadAsset(owner, repo, release.GetID(), path.Join(dir, file.Name()), accessToken)
			if err != nil {
				return err
			}
		}

		if prune {
			fmt.Println("checking for remote files to prune...")
			// Let's see if any remote assets need to be deleted.
			for _, asset := range assets {
				found := false
				for _, file := range files {
					if file.Name() == asset.GetName() {
						found = true
					}
				}
				if !found {
					fmt.Printf("pruning asset %s\n", asset.GetName())
					err = DeleteAsset(owner, repo, asset.GetID(), accessToken)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "github-release-sync: %s\n", err)
		os.Exit(1)
	}
}
