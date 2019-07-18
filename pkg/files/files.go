package files

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/VirtusLab/go-extended/pkg/errors"
)

// ErrExpectedStdin indicates that an stdin pipe was expected but not present
type ErrExpectedStdin struct {
	stack *errors.Stack
}

func (e *ErrExpectedStdin) Error() string {
	return "expected a pipe stdin"
}

// Format implements fmt.Formatter used by Sprint(f) or Fprint(f) etc.
func (e *ErrExpectedStdin) Format(s fmt.State, verb rune) {
	errors.FormatCauseAndStack(e, e.stack, s, verb)
}

// StackTrace returns a stack trace for this error
func (e *ErrExpectedStdin) StackTrace() errors.StackTrace {
	return e.stack.StackTrace()
}

// NewErrExpectedStdin creates a new ErrExpectedStdin
func NewErrExpectedStdin() *ErrExpectedStdin {
	return &ErrExpectedStdin{
		stack: errors.Callers(),
	}
}

// FileEntry contains file information
type FileEntry struct {
	Path      string
	Name      string
	Extension string
}

// ReadInput reads bytes from inputPath (if not empty) or stdin
func ReadInput(path string) ([]byte, error) {
	var inputFile *os.File
	if path == "" {
		stdinFileInfo, _ := os.Stdin.Stat()
		if (stdinFileInfo.Mode() & os.ModeNamedPipe) != 0 {
			inputFile = os.Stdin
		} else {
			return nil, NewErrExpectedStdin()
		}
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer func() { _ = f.Close() }()
		inputFile = f
	}
	fileContent, err := ioutil.ReadAll(inputFile)
	if err != nil {
		return nil, err
	}
	// golang adds a new line characters at the end of every line, not what we want here
	// note that we need to make sure the workaround is cross platform
	fileContent = bytes.TrimRight(fileContent, "\r\n")
	return fileContent, nil
}

// WriteOutput writes given bytes into outputPath (if not empty) or stdout
func WriteOutput(path string, contents []byte, perm os.FileMode) error {
	if path == "" {
		count, err := os.Stdout.Write(contents)
		if err == nil && count < len(contents) {
			return io.ErrShortWrite
		}
		if err != nil {
			return err
		}
	} else {
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, contents, perm)
		if err != nil {
			return errors.Wrapf(err, "can't write file: '%s'", path)
		}
	}
	return nil
}

// CheckNotEmptyAndExists returns an error if the given file does not exist exists or is empty
func CheckNotEmptyAndExists(file string) error {
	if len(file) == 0 {
		return errors.New("file path is empty")
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		return errors.New("file path does not exist")
	}

	if fileInfo.Size() == 0 {
		return errors.New("file is empty")
	}

	return nil
}

// ToAbsPath turns a relative path into an absolute path with the given root path, absolute paths are ignored
func ToAbsPath(path, root string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Join(root, path), nil
}

// Pwd returns the process working directory
func Pwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// DirTree returns files form all directories and subdirectories
func DirTree(input string) (entries []FileEntry, err error) {
	err = filepath.Walk(input, func(path string, info os.FileInfo, dirErr error) error {
		if dirErr != nil {
			return errors.Errorf("error '%v' on path '%s'", dirErr, path)
		}

		if !info.IsDir() {
			entry := FileEntry{
				Path:      filepath.Dir(path),
				Name:      info.Name(),
				Extension: filepath.Ext(path),
			}
			entries = append(entries, entry)
		}
		return nil
	})
	if err != nil {
		return entries, errors.Errorf("can't walk the directory tree '%s': %s", input, err)
	}

	return entries, nil
}

// TrimExtension returns file without given extensions
func TrimExtension(file FileEntry, extensions []string) (new FileEntry) {
	new = file
	for _, ext := range extensions {
		if file.Extension == ext {
			new.Name = strings.TrimSuffix(file.Name, ext)
			new.Extension = filepath.Ext(new.Name)
		}
	}
	return
}
