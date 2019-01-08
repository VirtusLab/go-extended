package check

import (
	"fmt"

	"github.com/VirtusLab/go-extended/pkg/matcher"
)

const (
	javaScriptIdentifierPattern = `^[a-zA-Z_$][a-zA-Z0-9_$]*$`
)

// InvalidJavaScriptIdentifier is used when the JavaScript identifier check fails
type InvalidJavaScriptIdentifier struct {
	text string
}

func (e *InvalidJavaScriptIdentifier) Error() string {
	return fmt.Sprintf("must be a valid JavaScript identifier, '%s' does not match pattern '%s'",
		e.text, javaScriptIdentifierPattern)
}

// IsValidJavaScriptIdentifier checks if the given string is a valid JavaScript/JSON identifier
func IsValidJavaScriptIdentifier(value string) error {
	if !matcher.Must(javaScriptIdentifierPattern).Match(value) {
		return &InvalidJavaScriptIdentifier{value}
	}
	return nil
}
