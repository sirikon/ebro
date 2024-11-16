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
		return nil, fmt.Errorf("parsing url %v: %w", importUrl, err)
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
		fragmentQuery, err := url.ParseQuery(strings.TrimPrefix(parsedUrl.Fragment, "?"))
		if err != nil {
			return nil, fmt.Errorf("parsing fragment of url %v: %w", importUrl, err)
		}

		if val, ok := fragmentQuery["ref"]; ok {
			setResultRef(plumbing.ReferenceName(val[0]))
		}
		if val, ok := fragmentQuery["path"]; ok {
			result.Subpath = val[0]
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
