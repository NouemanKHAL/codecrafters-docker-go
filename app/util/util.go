package util

import (
	"io"
	"os"
	"strings"
)

func JoinPath(parts ...string) string {
	path := strings.Join(parts, string(os.PathSeparator))
	return strings.Replace(path, "//", "/", -1)
}

func MoveFile(src, dst string) error {
	err := CopyFile(src, dst)
	if err != nil {
		return err
	}
	err = os.RemoveAll(src)
	return err
}

func CopyFile(src, dst string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstParentDir := GetParentDir(dst)
	err = os.MkdirAll(dstParentDir, 0750)
	if err != nil {
		return err
	}
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, srcStat.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func CreateFile(path string) (*os.File, error) {
	dirname := GetParentDir(path)
	err := os.MkdirAll(dirname, 0750)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0750)
	return f, err
}

func GetParentDir(filePath string) string {
	idx := strings.LastIndex(filePath, string(os.PathSeparator))
	return filePath[:idx]
}
