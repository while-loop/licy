package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"golang.org/x/oauth2"
	"context"
	"github.com/google/go-github/github"
	"log"
	"sync"
)

const (
	owner  = "while-loop"
	repo   = "licy"
	apiKey = "GITHUB_API_KEY"
)

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(apiKey)},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	ghClient := github.NewClient(tc)

	binDir := os.Args[1]
	tag := os.Args[2]
	releaseID, err := releaseID(ghClient, owner, repo, tag)
	if err != nil {
		log.Fatal(err)
	}

	bins, err := ioutil.ReadDir(binDir)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.Chdir(binDir); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, bin := range bins {
		wg.Add(1)
		go func(bin os.FileInfo) {
			defer wg.Done()
			fmt.Printf("publishing %v\n", bin.Name())
			f, err := os.Open(bin.Name())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			if err = publishBin(ghClient, owner, repo, releaseID, f); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}(bin)
	}

	wg.Wait()
}

func releaseID(client *github.Client, owner, repo, tag string) (int64, error) {
	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), owner, repo, tag)
	if err != nil {
		return 0, err
	}
	return *release.ID, nil
}

func publishBin(client *github.Client, owner, repo string, releaseID int64, bin *os.File) error {
	opt := &github.UploadOptions{Name: filepath.Base(bin.Name())}
	_, _, err := client.Repositories.UploadReleaseAsset(context.Background(), owner, repo, releaseID, opt, bin)
	return err
}
