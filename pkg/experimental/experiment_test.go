package experimental

import (
	"testing"

	"github.com/VirtusLab/go-extended/pkg/matcher"
)

func TestName(t *testing.T) {
	value := ""
	switch {
	default:
		t.Fatal("default")
	case matcher.Must(``).Match(value):

	}
}
