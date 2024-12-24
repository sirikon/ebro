package remote

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirikon/ebro/internal/logger"
)

type GitImport struct {
	Url     string
	Ref     plumbing.ReferenceName
	Path    string
	Subpath string
}

func ParseGitImport(importUrl string) (*GitImport, error) {
	if !strings.HasPrefix(importUrl, "git+") {
		return nil, nil
	}

	parsedUrl, err := url.Parse(strings.TrimPrefix(importUrl, "git+"))
	if err != nil {
		return nil, fmt.Errorf("parsing url: %w", err)
	}

	result := GitImport{
		Url:     parsedUrl.Scheme + "://" + parsedUrl.Host + parsedUrl.Path,
		Subpath: ".",
		Ref:     "",
		Path:    "",
	}

	setResultRef := func(ref plumbing.ReferenceName) {
		result.Ref = ref
		result.Path = path.Join(".ebro", "git", parsedUrl.Scheme, parsedUrl.Host+parsedUrl.Path, string(ref))
	}
	setResultRef(plumbing.Master)

	if parsedUrl.Fragment != "" {
		parsedFragmentUrl, err := url.Parse(parsedUrl.Fragment)
		if err != nil {
			return nil, fmt.Errorf("parsing fragment of url: %w", err)
		}
		parsedFragmentUrlQuery := parsedFragmentUrl.Query()
		for key, value := range parsedFragmentUrlQuery {
			if key == "ref" {
				setResultRef(plumbing.ReferenceName(value[0]))
				continue
			}
			return nil, fmt.Errorf("unknown query parameter in git import fragment: %v", key)
		}

		if parsedFragmentUrl.Path != "" {
			result.Subpath = parsedFragmentUrl.Path
		}
	}

	return &result, nil
}

func CloneGitImport(gi *GitImport) error {
	_, err := os.Stat(gi.Path)
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	logger.Notice("cloning " + gi.Url)
	_, err = git.PlainClone(gi.Path, false, &git.CloneOptions{
		URL:           gi.Url,
		ReferenceName: gi.Ref,
		SingleBranch:  true,
		Progress:      nil,
	})

	return err
}
