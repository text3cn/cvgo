package filekit

import (
	"cvgo/kit/strkit"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// 获取一个绝对路径所属目录
func Dir(absolutePath string) string {
	return filepath.Dir(absolutePath)
}

// 判断文件/文件夹是否存在
func PathExists(absolutePath string) (bool, error) {
	_, err := os.Stat(absolutePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 判断是否是文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断是否是文件而不是目录
func IsFile(path string) bool {
	return !IsDir(path)
}

// 根据文件名获取后缀
func GetSuffix(filename string) string {
	s := strkit.Explode(".", filename)
	last := s[len(s)-1]
	if last == filename {
		return ""
	}
	return last
}

// 创建目录
func MkDir(dir string, mode os.FileMode) {
	var err error
	if _, err = os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dir, mode); err != nil {
			fmt.Println("创建目录失败：" + dir)
		}
	}
}

// 扫描指定目录下所有文件及文件夹
// dir 指定要扫描的目录
// return：
// files 文件数组
// dirs  文件夹数组
func Scandir(dir string) (files []string, dirs []string) {
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			// 扫描到的 ./ (自身) 不要
			if path != dir {
				dirs = append(dirs, path)
			}
			return nil
		}
		files = append(files, path)
		return nil
	})
	return
}

// 创建文件，覆盖创建
func createFile(filepath string, content string) {
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write([]byte(content))
	}
}

// 读取文件内容
func readFile(filepath string) string {
	ret := ""
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		contentByte, _ := io.ReadAll(f)
		ret = string(contentByte)
	}
	return ret
}

// 上传文件保存到本地
func UploadFile(file *multipart.FileHeader) {
	// 输出文件信息
	fmt.Printf("Uploaded File: %+v\n", file.Filename)
	fmt.Printf("File Size: %+v\n", file.Size)
	fmt.Printf("MIME Header: %+v\n", file.Header)
}

// 判断是否图片类型
func IsImage(file *multipart.FileHeader) bool {
	// 获取文件的 MIME 类型
	mimeType := file.Header.Get("Content-Type")
	// 检查 MIME 类型是否是允许的图片类型
	return strings.HasPrefix(mimeType, "image/")
}

// 删除给定路径的文件
// filePath 文件的绝对路径
func DeleteFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

// 删除给定的文件或目录，如果不存在则直接返回 nil，如果存在则删除
func RemoveIfExists(path string) error {
	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 路径不存在，直接返回
		return nil
	} else if err != nil {
		// 其他错误，如权限问题
		return fmt.Errorf("error checking if path exists: %v", err)
	}
	// 路径存在，执行删除操作
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("error deleting path: %v", err)
	}
	return nil
}

// CreatePath 创建文件或目录
// 如果路径以 "/" 结尾，将被认为是目录；否则认为是文件。
func CreatePath(path string) error {
	// 判断路径是否是目录
	isDir := filepath.Ext(path) == ""

	if isDir {
		// 创建目录
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	} else {
		// 创建文件所在的目录
		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating parent directory: %v", err)
		}

		// 创建文件
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer file.Close()
	}

	return nil
}

// CopyFile 将 src 文件复制到 dst，如果目录不在会自动创建
func CopyFile(src, dst string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		fmt.Println("CopyFile 打开源文件失败", err)
		return err
	}
	defer sourceFile.Close()

	// 确保目标目录存在
	destDir := filepath.Dir(dst)
	if err := EnsureDirExists(destDir); err != nil {
		fmt.Println("EnsureDirExists Error:", err)
		return err
	}

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("CopyFile 创建目标文件失败", err)
		return err
	}
	defer destFile.Close()

	// 将源文件内容复制到目标文件
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		fmt.Println("CopyFile 将源文件内容复制到目标文件失败", err)
		return err
	}

	// 确保写入完成
	err = destFile.Sync()
	if err != nil {
		return err
	}

	// 复制文件权限
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, fileInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

// EnsureDirExists 确保目标目录存在，如果不存在则创建它
func EnsureDirExists(dir string) error {
	// 检查目录是否存在，不存在则创建
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}
	return nil
}

// CopyFiles 将 srcDir 目录下的所有文件复制到 destDir 目录下
func CopyFiles(srcDir, dstDir string) error {
	// 检查目标目录是否存在，不存在则创建
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err = os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create destination directory: %v", err)
		}
	}

	// 遍历源目录中的所有文件
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("遍历出错:", err)
			return err
		}

		// 忽略目录，只复制文件
		if !info.IsDir() {
			// 生成目标文件的路径
			relPath, err := filepath.Rel(srcDir, path)
			if err != nil {
				fmt.Printf("Failed to get relative path: %+v\n", err)
				return err
			}
			destPath := filepath.Join(dstDir, relPath)

			// 创建目标文件所在的目录结构
			destDirPath := filepath.Dir(destPath)
			if _, err := os.Stat(destDirPath); os.IsNotExist(err) {
				err = os.MkdirAll(destDirPath, os.ModePerm)
				if err != nil {
					fmt.Printf("Failed to create destination directory: %+v\n", err)
					return err
				}
			}

			// 复制文件
			err = CopyFile(path, destPath)
			if err != nil {
				fmt.Printf("Failed to copy file: %+v\n", err)
				return err
			}
		}
		return nil
	})

	return err
}

// 移动目录，目标目录不存在会创建
func MoveDir(src, dst string) error {
	// 检查源路径是否存在
	_, err := os.Stat(src)
	if os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", src)
	}
	// 确保目标目录存在，不存在就创建
	EnsureDirExists(dst)
	dst = filepath.Join(dst, filepath.Base(src))
	// 移动
	err = os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("failed to move %s to %s: %v", src, dst, err)
	}
	return nil
}

// MoveFiles 将 srcDir 目录下的所有文件移动到 destDir 目录下
func MoveFiles(srcDir, destDir string) error {
	// 检查 srcDir 是否存在
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", srcDir)
	}

	// 检查 destDir 是否存在，如果不存在则创建它
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err := os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating destination directory: %v", err)
		}
	}

	// 读取 srcDir 目录下的所有文件
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("error reading source directory: %v", err)
	}

	// 遍历文件列表并移动它们到目标目录
	for _, file := range files {
		// 构建源文件路径和目标文件路径
		srcFilePath := filepath.Join(srcDir, file.Name())
		destFilePath := filepath.Join(destDir, file.Name())

		// 移动文件
		err := os.Rename(srcFilePath, destFilePath)
		if err != nil {
			return fmt.Errorf("error moving file %s: %v", file.Name(), err)
		}
	}
	return nil
}

// Rename 重命名目录或文件
func Rename(oldDir, newDir string) error {
	// 使用 os.Rename 进行重命名操作
	err := os.Rename(oldDir, newDir)
	if err != nil {
		return fmt.Errorf("failed to rename directory from %s to %s: %v", oldDir, newDir, err)
	}
	return nil
}
