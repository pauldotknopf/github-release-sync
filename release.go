package main

import (
	"context"
	"os"
	"path"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GetReleasesTag GetReleasesTag
func GetReleasesTag(owner, repo, tag string) (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	repoRelease, _, err := client.Repositories.GetReleaseByTag(context.Background(), owner, repo, tag)
	if err != nil {
		return nil, err
	}
	return repoRelease, nil
}

// GetReleaseAssets GetReleaseAssets
func GetReleaseAssets(owner, repo string, releaseID int64) ([]*github.ReleaseAsset, error) {
	client := github.NewClient(nil)
	assets, _, err := client.Repositories.ListReleaseAssets(context.Background(),
		owner,
		repo,
		releaseID,
		nil)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// DeleteAsset DeleteAsset
func DeleteAsset(owner, repo string, assetID int64, accessToken string) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	_, err := client.Repositories.DeleteReleaseAsset(context.Background(), owner, repo, assetID)
	return err
}

// UploadAsset UploadAsset
func UploadAsset(owner, repo string, releaseID int64, filePath, accessToken string) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_, _, err = client.Repositories.UploadReleaseAsset(context.Background(),
		owner,
		repo,
		releaseID,
		&github.UploadOptions{
			Name: path.Base(filePath),
		},
		file)
	return err
}
