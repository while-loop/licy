package main

import (
	"github.com/sajari/fuzzy"
	"fmt"
	"strings"
	"io"
)

func suggestMispell(writer io.Writer, license string) {
	model := fuzzy.NewModel()
	model.SetThreshold(1)
	model.SetDepth(4)
	lics, err := getSuggestions()
	if err != nil {
		fmt.Fprintln(writer, err)
		return
	}
	model.Train(lics)
	suggs := model.Suggestions(license, false)
	if len(suggs) > 0 {
		fmt.Fprintf(writer, "\"%s\". Did you mean \"%s\"?\n", license, suggs[0])
	}
}

func getSuggestions() ([]string, error) {
	lics := make([]string, 0)
	licenses, err := GetLicenses()
	if err != nil {
		return nil, err
	}

	for _, l := range licenses {
		lics = append(lics, l.SpdxID, strings.ToLower(l.SpdxID))
	}

	return lics, nil
}
