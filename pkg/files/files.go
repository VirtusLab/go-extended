package files

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	// ErrExpectedStdin indicates that an stdin pipe was expected but not present
	ErrExpectedStdin = errors.New("expected a pipe stdin")
)

// ReadInput reads bytes from inputPath (if not empty) or stdin
func ReadInput(inputPath string) ([]byte, error) {
	var inputFile *os.File
	if inputPath == "" {
		stdinFileInfo, _ := os.Stdin.Stat()
		if (stdinFileInfo.Mode() & os.ModeNamedPipe) != 0 {
			inputFile = os.Stdin
		} else {
			return nil, ErrExpectedStdin
		}
	} else {
		f, err := os.Open(inputPath)
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
	return fileContent, nil
}

// WriteOutput writes given bytes into outputPath (if not empty) or stdout
func WriteOutput(outputPath string, outputContent []byte, perm os.FileMode) error {
	if outputPath == "" {
		count, err := os.Stdout.Write(outputContent)
		if err == nil && count < len(outputContent) {
			return io.ErrShortWrite
		}
		if err != nil {
			return err
		}
	} else {
		err := ioutil.WriteFile(outputPath, outputContent, perm)
		if err != nil {
			return err
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
