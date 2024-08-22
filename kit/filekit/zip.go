package filekit

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipDirectory 将 srcDir 目录打包压缩，并输出到 dstZipFile
// e.g. ZipDirectory(dist, dist+"dist.zip")
func ZipDirectory(srcDir, dstZipFile string) error {
	// 创建目标 ZIP 文件
	zipFile, err := os.Create(dstZipFile)
	if err != nil {
		return fmt.Errorf("failed to add zip file: %v", err)
	}
	defer zipFile.Close()

	// 创建 ZIP Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历源目录并将文件添加到 ZIP
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是目标 ZIP 文件，避免死循环
		if path == dstZipFile {
			return nil
		}

		// 获取相对路径并将其标准化为 Unix 路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		relPath = strings.ReplaceAll(relPath, string(filepath.Separator), "/")

		// 如果是目录，添加到 ZIP 但不压缩内容
		if info.IsDir() {
			_, err = zipWriter.Create(relPath + "/")
			if err != nil {
				return err
			}
			return nil
		}

		// 如果是文件，添加到 ZIP 并压缩内容
		zipFileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// 打开源文件
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 将源文件内容写入 ZIP 文件中
		_, err = io.Copy(zipFileWriter, srcFile)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to zip directory: %v", err)
	}

	return nil
}
