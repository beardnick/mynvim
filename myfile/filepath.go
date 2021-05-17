package myfile

import (
	"errors"
	"os"
	"path/filepath"
)

// 不存在就创建文件夹
func EnsureDir(path string) (err error) {
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, os.ModePerm)
		return
	}
	return
}


// 不存在就创建文件
func EnsureFile(path string) (err error) {
	if !FileExists(path) {
		return
	}
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}
	_, err = os.Create(path)
	return
}

func FileExists(path string)  bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrNotExist)
}