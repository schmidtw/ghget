package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/schmidtw/ghget"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Fprintf(os.Stderr, "usage: %s <release-asset-url>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  e.g. https://github.com/OWNER/REPO/releases/download/TAG/ASSET")
		fmt.Fprintln(os.Stderr, "  set GITHUB_TOKEN in the environment for private repos")
		os.Exit(2)
	}

	rawURL := os.Args[1]
	name, err := assetName(rawURL)
	if err != nil {
		die(err)
	}

	fmt.Printf("downloading %s\n", name)
	if err := ghget.Download(rawURL, name, os.Getenv("GITHUB_TOKEN")); err != nil {
		die(err)
	}
}

func assetName(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	name := path.Base(strings.Trim(u.Path, "/"))
	if name == "" || name == "." || name == "/" {
		return "", fmt.Errorf("could not determine asset name from URL: %s", raw)
	}
	return name, nil
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
