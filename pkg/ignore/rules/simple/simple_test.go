package simple

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/VirtusLab/go-extended/pkg/ignore"

	"github.com/stretchr/testify/assert"
)

type MockFileInfo struct {
	isDir bool
}

func (f MockFileInfo) Name() string {
	panic("implement me")
}

func (f MockFileInfo) Size() int64 {
	panic("implement me")
}

func (f MockFileInfo) Mode() os.FileMode {
	panic("implement me")
}

func (f MockFileInfo) ModTime() time.Time {
	panic("implement me")
}

func (f MockFileInfo) Sys() interface{} {
	panic("implement me")
}

func (f MockFileInfo) IsDir() bool {
	return f.isDir
}

func TestMatch(t *testing.T) {
	data := []struct {
		rules []string
		path  string
		isDir bool
		want  bool
	}{
		{[]string{"a.txt"}, "a.txt", false, true},
		{[]string{"*.txt"}, "a.txt", false, true},
		{[]string{"dir/a.txt"}, "dir/a.txt", false, true},
		{[]string{"dir/*.txt"}, "dir/a.txt", false, true},
		{[]string{"dir2/a.txt"}, "dir1/dir2/a.txt", false, true},
		{[]string{"dir3/a.txt"}, "dir1/dir2/dir3/a.txt", false, true},
		{[]string{"a.txt"}, "dir/a.txt", false, true},
		{[]string{"*.txt"}, "dir/a.txt", false, true},
		{[]string{"a.txt"}, "dir1/dir2/a.txt", false, true},
		{[]string{"dir2/a.txt"}, "dir1/dir2/a.txt", false, true},
		{[]string{"dir"}, "dir", true, true},
		{[]string{"dir/"}, "dir", true, true},
		{[]string{"dir/"}, "dir", false, false},
		{[]string{"dir1/dir2/"}, "dir1/dir2", true, true},
		{[]string{"/a.txt"}, "a.txt", false, true},
		{[]string{"/dir/a.txt"}, "dir/a.txt", false, true},
		{[]string{"/dir1/a.txt"}, "dir/dir1/a.txt", false, false},
		{[]string{"/a.txt"}, "dir/a.txt", false, false},
		{[]string{"a.txt", "b.txt"}, "dir/b.txt", false, true},
		{[]string{"*.txt", "!b.txt"}, "dir/b.txt", false, false},
		{[]string{"dir/*.txt", "!dir/b.txt"}, "dir/b.txt", false, false},
		{[]string{"dir/*.txt", "!/b.txt"}, "dir/b.txt", false, true},
		// blacklisting
		{[]string{".*", "!.want*"}, ".wanted", false, false},
		{[]string{".*", "!.want*"}, ".notsomuch", false, true},
		{[]string{".*!@#$%^&*()_+{]\\//[}"}, ".", false, false},
	}

	for i, sample := range data {
		t.Run(fmt.Sprintf("[%d] %s %s", i, sample.rules, sample.path), func(t *testing.T) {
			rules, err := Parse(strings.NewReader(strings.Join(sample.rules, "\n")))
			assert.NoError(t, err)

			got, err := ignore.New(".", rules).Match(sample.path, MockFileInfo{sample.isDir})

			assert.NoError(t, err)
			assert.EqualValues(t, sample.want, got, fmt.Sprintf("%+v", sample))
		})
	}
}

func TestCutLastN(t *testing.T) {
	data := []struct {
		path string
		n int
		want string
	}{
		{"", 0, ""},
		{"\\//\\//", 0, ""},
		{"/one/two/three/four", -1, ""},
		{"/one/two/three/four", 0, ""},
		{"/one/two/three/four", 1, "four"},
		{"one/two/three/four", 2, "three/four"},
		{"one/two/three/four", 3, "two/three/four"},
		{"one/two/three/four", 4, "one/two/three/four"},
		{"one/two/three/four", 666, "one/two/three/four"},
	}

	for i, sample := range data {
		t.Run(fmt.Sprintf("[%d] %s %d %s", i, sample.path, sample.n, sample.want), func(t *testing.T) {
			got := cutLastN(sample.path, sample.n)
			assert.Equal(t, sample.want, got)
		})
	}
}
