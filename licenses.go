package main

import (
	"gopkg.in/yaml.v1"
	"fmt"
	"os"
	"sort"
	"bytes"
	"strings"
	"strconv"
	"github.com/pkg/errors"
	"io"
)

const (
	licenseDir = "_licenses/"
)

var (
	ErrNotFound = errors.New("license not found")
)

type License struct {
	Conditions  []string `yaml:"conditions"`
	Description string   `yaml:"description"`
	How         string   `yaml:"how"`
	Limitations []string `yaml:"limitations"`
	Permissions []string `yaml:"permissions"`
	Source      string   `yaml:"source"`
	SpdxID      string   `yaml:"spdx-id"`
	Title       string   `yaml:"title"`
	Using       string   `yaml:"using"`
	Body        string
}

type Licenses []License

func (l Licenses) Len() int {
	return len(l)
}

func (l Licenses) Less(i, j int) bool {
	return l[i].Title < l[j].Title
}

func (l Licenses) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func GetLicenses() ([]License, error) {
	licenses := make([]License, 0)

	for _, l := range AssetNames() {

		bs, err := Asset(l)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to get %s: %v\n", l, err)
			continue
		}
		lic, err := parseLicense(bs)
		licenses = append(licenses, lic)
	}

	sort.Sort(Licenses(licenses))
	return licenses, nil
}

func GetLicense(licenseName string) (License, error) {
	bs, err := Asset(licenseDir + licenseName + ".txt")
	if err != nil {
		return License{}, ErrNotFound
	}

	return parseLicense(bs)
}

func (l *License) FillBody(org string, year int) {
	yearTempl, nameTempl := getTemplate(l.SpdxID)
	if yearTempl != "" {
		l.Body = strings.Replace(l.Body, yearTempl, strconv.Itoa(year), -1)
	}

	if nameTempl != "" {
		l.Body = strings.Replace(l.Body, nameTempl, org, -1)
	}
}

func getTemplate(spdxID string) (string, string) {
	switch strings.ToLower(spdxID) {
	case "apache-2.0":
		return "[yyyy]", "[name of copyright owner]"
	case "bsd-2-clause", "bsd-3-clause", "mit", "bsd-3-clause-clear":
		return "[year]", "[fullname]"
	case "lgpl-2.1", "gpl-3.0", "gpl-2.0":
		return "<year>", "<name of author>"
	}

	return "", ""
}

func (l *License) FillProject(project string) {
	l.Body = strings.Replace(l.Body, "[project]", project, -1)
}

func (l *License) Save(w io.Writer) (int, error) {
	return io.WriteString(w, l.Body)
}

func parseLicense(bs []byte) (License, error) {
	var lic License

	bodyIdx := bytes.LastIndex(bs, []byte("---"))
	if err := yaml.Unmarshal(bs[:bodyIdx], &lic); err != nil {
		return License{}, errors.Wrap(err, "Unable to unmarshal license")
	}

	lic.Body = strings.Trim(string(bs[bodyIdx+3:]), "\r\n")
	return lic, nil
}
