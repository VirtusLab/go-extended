package ignore

import (
	"os"
	"path/filepath"
)

type Rules interface {
	Match(path string, info os.FileInfo) bool
}

type Matcher interface {
	Match(path string, info os.FileInfo) (bool, error)
}

type matcher struct {
	root  string
	rules Rules
}

func New(root string, rules Rules) Matcher {
	return &matcher{
		root:  root,
		rules: rules,
	}
}

func (m matcher) Match(path string, info os.FileInfo) (bool, error) {
	relativePath, err := filepath.Rel(m.root, path)
	if err != nil {
		return false, err
	}
	return m.rules.Match(relativePath, info), nil
}
