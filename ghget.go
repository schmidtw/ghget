// Package ghget downloads release assets from GitHub.
package ghget

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type asset struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
	URL  string `json:"url"`
}

type release struct {
	TagName string  `json:"tag_name"`
	Assets  []asset `json:"assets"`
}

// Download fetches the GitHub release asset at rawURL and writes it to dstPath.
// When token is non-empty it is used for authentication (required for private repos).
// rawURL must be of the form https://github.com/OWNER/REPO/releases/download/TAG/ASSET.
func Download(rawURL, dstPath, token string) error {
	owner, repo, tag, name, err := parseReleaseURL(rawURL)
	if err != nil {
		return err
	}

	downloadURL := rawURL
	if token != "" {
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)
		rel, err := fetchRelease(apiURL, token)
		if err != nil {
			return err
		}
		var found *asset
		for i := range rel.Assets {
			if rel.Assets[i].Name == name {
				found = &rel.Assets[i]
				break
			}
		}
		if found == nil {
			return fmt.Errorf("asset %q not found in release %s", name, rel.TagName)
		}
		downloadURL = found.URL
	}

	return fetchToFile(downloadURL, dstPath, token)
}

func parseReleaseURL(raw string) (owner, repo, tag, name string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", "", "", err
	}
	if u.Host != "github.com" {
		return "", "", "", "", fmt.Errorf("expected github.com URL, got %q", u.Host)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 6 || parts[2] != "releases" || parts[3] != "download" {
		return "", "", "", "", fmt.Errorf("not a release asset URL: %s", raw)
	}
	return parts[0], parts[1], parts[4], parts[5], nil
}

func fetchRelease(url, token string) (*release, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api %d: %s", resp.StatusCode, body)
	}
	var r release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func fetchToFile(url, dst, token string) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/octet-stream")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("download %s: status %d", url, resp.StatusCode)
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
