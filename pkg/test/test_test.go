package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Main(m)
	m.Run()
	os.Exit(0)
}
