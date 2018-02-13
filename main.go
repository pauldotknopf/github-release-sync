package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/urfave/cli"
)

func main() {
	var accessToken string
	var owner string
	var repo string
	var tag string
	var dir string

	app := cli.NewApp()
	app.Name = "darch"
	app.Usage = "A tool used to build, boot and share stateless Arch images."
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "access-token",
			Usage:       "the github access token to use the api",
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
		cli.StringFlag{
			Name:        "dir",
			Usage:       "the directory to sync to the github release",
			Value:       ".",
			Destination: &dir,
		},
	}
	app.Action = func(clicontext *cli.Context) error {
		if len(accessToken) == 0 {
			return fmt.Errorf("access token is required")
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
			found := false
			for _, asset := range assets {
				if asset.GetName() == file.Name() {
					found = true
				}
			}
			if found {
				fmt.Printf("%s already uploaded, skipping...\n", file.Name())
				continue
			}
			fmt.Printf("uplading %s\n", file.Name())
			err = UploadAsset(owner, repo, release.GetID(), path.Join(dir, file.Name()), accessToken)
			if err != nil {
				return err
			}
		}

		fmt.Println(release.GetTagName())

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "github-release-sync: %s\n", err)
		os.Exit(1)
	}
}
