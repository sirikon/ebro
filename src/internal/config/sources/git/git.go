package git

// import (
// 	"errors"
// 	"fmt"
// 	"net/url"
// 	"os"
// 	"path"
// 	"strings"

// 	"github.com/go-git/go-git/v5"
// 	"github.com/go-git/go-git/v5/plumbing"
// 	"github.com/sirikon/ebro/internal/logger"
// )

// type GitSource struct{}

// func (gs GitSource) Match(from string) (bool, error) {
// 	return isGitImport(from), nil
// }

// func (gs GitSource) Resolve(_ string, from string) (string, error) {
// 	parsedGitImport, err := parseGitImport(from)
// 	if err != nil {
// 		return "", fmt.Errorf("parsing git import: %w", err)
// 	}

// 	if parsedGitImport == nil {
// 		return "", fmt.Errorf("parsing git import returned nothing")
// 	}

// 	err = cloneGitImport(parsedGitImport)
// 	if err != nil {
// 		return "", fmt.Errorf("cloning git import: %w", err)
// 	}

// 	return path.Join(parsedGitImport.path, parsedGitImport.subpath), nil
// }

// type gitImport struct {
// 	url     string
// 	ref     plumbing.ReferenceName
// 	path    string
// 	subpath string
// }

// func isGitImport(from string) bool {
// 	return strings.HasPrefix(from, "git+")
// }

// func parseGitImport(importUrl string) (*gitImport, error) {
// 	parsedUrl, err := url.Parse(strings.TrimPrefix(importUrl, "git+"))
// 	if err != nil {
// 		return nil, fmt.Errorf("parsing url: %w", err)
// 	}

// 	result := gitImport{
// 		url:     parsedUrl.Scheme + "://" + parsedUrl.Host + parsedUrl.Path,
// 		subpath: ".",
// 		ref:     "",
// 		path:    "",
// 	}

// 	setResultRef := func(ref plumbing.ReferenceName) {
// 		result.ref = ref
// 		result.path = path.Join(".ebro", "git", parsedUrl.Scheme, parsedUrl.Host+parsedUrl.Path, string(ref))
// 	}
// 	setResultRef(plumbing.Master)

// 	if parsedUrl.Fragment != "" {
// 		parsedFragmentUrl, err := url.Parse(parsedUrl.Fragment)
// 		if err != nil {
// 			return nil, fmt.Errorf("parsing fragment of url: %w", err)
// 		}
// 		parsedFragmentUrlQuery := parsedFragmentUrl.Query()
// 		for key, value := range parsedFragmentUrlQuery {
// 			if key == "ref" {
// 				setResultRef(plumbing.ReferenceName(value[0]))
// 				continue
// 			}
// 			return nil, fmt.Errorf("unknown query parameter in git import fragment: %v", key)
// 		}

// 		if parsedFragmentUrl.Path != "" {
// 			result.subpath = parsedFragmentUrl.Path
// 		}
// 	}

// 	return &result, nil
// }

// func cloneGitImport(gi *gitImport) error {
// 	_, err := os.Stat(gi.path)
// 	if err == nil {
// 		return nil
// 	}

// 	if !errors.Is(err, os.ErrNotExist) {
// 		return err
// 	}

// 	logger.Notice("cloning " + gi.url)
// 	_, err = git.PlainClone(gi.path, false, &git.CloneOptions{
// 		URL:           gi.url,
// 		ReferenceName: gi.ref,
// 		SingleBranch:  true,
// 		Progress:      nil,
// 	})

// 	return err
// }
