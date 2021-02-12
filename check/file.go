package check

import (
	"github.com/go-git/go-git/v5"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func RepoExists(path string) bool {
	_, err := git.PlainOpen(path)
	ok := true
	if err, ok = err.(git.NoMatchingRefSpecError); ok {
		return false
	}
	return true
}
