package config

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

type gitImport struct {
	url     string
	ref     plumbing.ReferenceName
	path    string
	subpath string
}

func parseGitImport(importUrl string) (*gitImport, error) {
	if !strings.HasPrefix(importUrl, "git+") {
		return nil, nil
	}

	parsedUrl, err := url.Parse(strings.TrimPrefix(importUrl, "git+"))
	if err != nil {
		return nil, fmt.Errorf("parsing url %v: %w", importUrl, err)
	}

	result := gitImport{
		url:     parsedUrl.Scheme + "://" + parsedUrl.Host + parsedUrl.Path,
		subpath: ".",
		ref:     "",
		path:    "",
	}

	setResultRef := func(ref plumbing.ReferenceName) {
		result.ref = ref
		result.path = path.Join(".ebro", "git", parsedUrl.Scheme, parsedUrl.Host+parsedUrl.Path, string(ref))
	}
	setResultRef(plumbing.Master)

	if parsedUrl.Fragment != "" {
		fragmentQuery, err := url.ParseQuery(strings.TrimPrefix(parsedUrl.Fragment, "?"))
		if err != nil {
			return nil, fmt.Errorf("parsing fragment of url %v: %w", importUrl, err)
		}

		if val, ok := fragmentQuery["ref"]; ok {
			setResultRef(plumbing.ReferenceName(val[0]))
		}
		if val, ok := fragmentQuery["path"]; ok {
			result.subpath = val[0]
		}
	}

	return &result, nil
}

func cloneGitImport(gi *gitImport) error {
	_, err := os.Stat(gi.path)
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	logger.Notice("cloning " + gi.url)
	_, err = git.PlainClone(gi.path, false, &git.CloneOptions{
		URL:           gi.url,
		ReferenceName: gi.ref,
		SingleBranch:  true,
		Progress:      os.Stderr,
	})

	return err
}
