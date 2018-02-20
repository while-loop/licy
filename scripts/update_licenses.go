package main

import (
	"github.com/google/go-github/github"
	"context"
	"fmt"
	"os"
	"io"
	"sync"
	"strings"
	"golang.org/x/oauth2"
)

const (
	owner  = "github"
	repo   = "choosealicense.com"
	path   = "_licenses/"
	apiKey = "GITHUB_API_KEY"
)

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(apiKey)},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	ghClient := github.NewClient(tc)

	_, files, _, err := ghClient.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		perr("Unable to reach choosealicense. " + err.Error())
	}

	if err = os.Mkdir(path, 0777); err != nil && !strings.Contains(err.Error(), "exists") {
		perr("Unable to create license dir: " + err.Error())
	}

	var wg sync.WaitGroup

	for _, name := range files {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			filename := path + name
			rc, err := ghClient.Repositories.DownloadContents(context.Background(), owner, repo, filename, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to download %s: %v\n", filename, err)
				return
			}
			defer rc.Close()

			f, err := os.Create(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create file %s: %v\n", filename, err)
				return
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to copy file contents to %s: %v\n", filename, err)
				return
			}
		}(name.GetName())
	}

	wg.Wait()
}

func perr(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
