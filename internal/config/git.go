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

type gitReference struct {
	cloneUrl  string
	clonePath string
	branch    string
	subPath   string
}

func parseGitReference(refUrl string) (*gitReference, error) {
	if !strings.HasPrefix(refUrl, "git+") {
		return nil, nil
	}
	refUrl = strings.TrimPrefix(refUrl, "git+")
	parsedUrl, err := url.Parse(refUrl)
	if err != nil {
		return nil, fmt.Errorf("parsing url %v: %w", refUrl, err)
	}
	branch := "master"
	subPath := "."
	if parsedUrl.Fragment != "" {
		values, err := url.ParseQuery(strings.TrimPrefix(parsedUrl.Fragment, "?"))
		if err != nil {
			return nil, fmt.Errorf("parsing fragment of url %v: %w", refUrl, err)
		}
		if val, ok := values["path"]; ok {
			subPath = val[0]
		}
		if val, ok := values["branch"]; ok {
			branch = val[0]
		}
	}

	return &gitReference{
		cloneUrl:  parsedUrl.Scheme + "://" + parsedUrl.Host + parsedUrl.Path,
		clonePath: path.Join(".ebro", "git", parsedUrl.Scheme, parsedUrl.Host+parsedUrl.Path, branch),
		branch:    branch,
		subPath:   subPath,
	}, nil
}

func cloneGitReference(ref *gitReference) error {
	_, err := os.Stat(ref.clonePath)
	if err == nil {
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	logger.Notice("cloning " + ref.cloneUrl)
	_, err = git.PlainClone(ref.clonePath, false, &git.CloneOptions{
		URL:           ref.cloneUrl,
		ReferenceName: plumbing.ReferenceName(ref.branch),
		SingleBranch:  true,
		Progress:      os.Stderr,
	})
	return err
}
