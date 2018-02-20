package main

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

var (
	owner        = "licenses"
	repo         = "license-templates"
	templatePath = "templates"
)

func downloadLicense(client *http.Client, license string) (io.ReadCloser, error) {
	ghClient := github.NewClient(client)
	if l, exists := licenseMap[license]; exists {
		license = l
	}
	if !strings.HasSuffix(license, ".txt") {
		license += ".txt"
	}

	templatePath += "/" + license

	rc, err := ghClient.Repositories.DownloadContents(context.Background(), owner, repo, templatePath,
		&github.RepositoryContentGetOptions{Ref: "aa0399cd31350a2692d0f51f651fc2fd3d0a5dab"})
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to download license file: %s", license)
	}

	return rc, nil
}

var (
	licenseMap = map[string]string{
		"apache2": "apache",
	}
)
