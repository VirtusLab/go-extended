package simple

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/VirtusLab/go-extended/pkg/ignore"
)

type rules struct {
	accepted ruleSet
	ignored  ruleSet
}

type rule struct {
	pattern       string
	hasRootPrefix bool
	hasDirSuffix  bool
	pathDepth     int
}

type ruleSet []rule

func parseRule(pattern string) (rule, error) {
	hasRootPrefix := strings.HasPrefix(pattern, "/")
	hasDirSuffix := strings.HasSuffix(pattern, "/")

	var pathDepth int
	if !hasRootPrefix {
		pathDepth = strings.Count(pattern, "/")
	}

	trimmedPattern := strings.Trim(pattern, "/")

	return rule{
		pattern:       trimmedPattern,
		hasRootPrefix: hasRootPrefix,
		hasDirSuffix:  hasDirSuffix,
		pathDepth:     pathDepth,
	}, validateMatch(trimmedPattern)
}

func Parse(reader io.Reader) (ignore.Rules, error) {
	var accepted ruleSet
	var ignored ruleSet

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// ignore empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "!") {
			rule, err := parseRule(strings.TrimPrefix(line, "!"))
			if err != nil {
				return nil, fmt.Errorf("rule '%s' is invalid: %s", line, err.Error())
			}
			accepted = append(accepted, rule)
		} else {
			rule, err := parseRule(line)
			if err != nil {
				return nil, fmt.Errorf("rule '%s' is invalid: %s", line, err.Error())
			}
			ignored = append(ignored, rule)
		}
	}
	return &rules{
		ruleSet(accepted),
		ruleSet(ignored),
	}, nil
}

func (r rules) Match(path string, info os.FileInfo) bool {
	if r.accepted.match(path, info.IsDir()) {
		return false
	}
	return r.ignored.match(path, info.IsDir())
}

func (rs ruleSet) match(path string, isDir bool) bool {
	for _, r := range rs {
		match := r.match(path, isDir)
		log.Printf("%+v | %s %v -> %v", r, path, isDir, match)
		if match {
			return true
		}
	}
	return false
}

func (r rule) match(path string, isDir bool) bool {
	if r.hasDirSuffix && !isDir {
		log.Printf("hasDirSuffix && !isDir")
		return false
	}

	targetPath := path
	if !r.hasRootPrefix { // relative path pattern
	    log.Printf("relative path pattern")
		targetPath = cutLastN(path, r.pathDepth+1)
	}
	log.Printf("targetPath = %v", targetPath)

	return mustMatch(r.pattern, targetPath)
}

func cutLastN(path string, n int) string {
	parts := strings.Split(path, "/")
	count := len(parts)
	if n < 0 {
		n = 0
	}
	if count > n {
		return strings.Join(parts[count-n:count], "/")
	} else {
		return path
	}
}

func mustMatch(pattern, name string) bool {
	matched, err := filepath.Match(pattern, name)
	if err != nil {
		panic("Unexpected error: " + err.Error())
	}
	return matched
}

func validateMatch(pattern string) error {
	_, err := filepath.Match(pattern, "")
	return err
}