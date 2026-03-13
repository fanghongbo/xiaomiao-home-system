package utils

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"mime/multipart"
	"os"
)

// GetMultipartFileMd5 获取上传文件md5
func GetMultipartFileMd5(fs multipart.File) (string, error) {
	var (
		buf     []byte
		hashVal hash.Hash
		hashStr string
		err     error
	)

	buf = make([]byte, 1024)
	hashVal = md5.New()

	for {
		var count int

		if count, err = fs.Read(buf); err != nil && err != io.EOF {
			return "", err
		}

		if count == 0 {
			break
		}

		hashVal.Write(buf[:count])
	}

	hashStr = fmt.Sprintf("%x", string(hashVal.Sum([]byte(""))))

	return hashStr, nil
}

// GetFileInfo 查询文件信息
func GetFileInfo(path string) (os.FileInfo, error) {
	var (
		fs  os.FileInfo
		err error
	)

	if fs, err = os.Stat(path); err != nil {
		return nil, err
	}

	return fs, nil
}

// ReadDir 读取目录内容
func ReadDir(path string) ([]os.FileInfo, error) {
	var (
		file    *os.File
		entries []os.FileInfo
		err     error
	)

	if file, err = os.Open(path); err != nil {
		return nil, err
	}

	defer file.Close()

	// 读取目录中的文件和子目录
	if entries, err = file.Readdir(-1); err != nil {
		return nil, err
	}

	return entries, nil
}
