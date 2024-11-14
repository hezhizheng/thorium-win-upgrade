package helper

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 比较版本号
// https://leetcode-cn.com/problems/compare-version-numbers/solution/golangshi-xian-by-he-qing-ping/
func CompareVersion(version1 string, version2 string) int {
	versionA := strings.Split(version1, ".")
	versionB := strings.Split(version2, ".")

	for i := len(versionA); i < 4; i++ {
		versionA = append(versionA, "0")
	}
	for i := len(versionB); i < 4; i++ {
		versionB = append(versionB, "0")
	}
	for i := 0; i < 4; i++ {
		version1, _ := strconv.Atoi(versionA[i])
		version2, _ := strconv.Atoi(versionB[i])
		if version1 == version2 {
			continue
		} else if version1 > version2 {
			return 1
		} else {
			return -1
		}
	}
	return 0
}

// Unzip 解压ZIP文件到指定目录
func Unzip(zipPath, destDir string) error {
	// 打开ZIP文件
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// 遍历ZIP文件中的文件
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// 构建解压后的文件路径
		path := filepath.Join(destDir, f.Name)

		// 如果是目录，则创建目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			// 创建文件
			os.MkdirAll(filepath.Dir(path), f.Mode())
			of, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer of.Close()

			// 将内容写入文件
			_, err = io.Copy(of, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
