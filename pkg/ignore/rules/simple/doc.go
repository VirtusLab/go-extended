/*
Package simple provides default rules for ignore.Matcher

Simple rules (ignore/rules/simple.Parse):
- Support absolute path (/path/to/ignore)
- Support relative path (path/to/ignore)
- Support accept pattern (!path/to/accept)
- Support directory pattern (path/to/directory/)
- Support glob pattern (path/to/*.txt)
- NO support for double-star glob pattern (uses filepath.Match) (path\/**\/file)
*/
package simple
