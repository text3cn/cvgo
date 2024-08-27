package tpl

import (
	"embed"
	"github.com/textthree/cvgokit/filekit"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed work work/**/* work/**/*.gitkeep work/.gitignore
var WorkDirTpl embed.FS

//go:embed module module/internal/**/*.gitkeep
var ModuleTpl embed.FS

//go:embed enable/gorm/entity/mysql/base.go
var EntityBase embed.FS

//go:embed enable/gorm/entity/autoMigrate.go
var AutoMigrate embed.FS

//go:embed enable/gorm/entityregistry/entityRegistry.go
var EntityRegistry embed.FS

//go:embed enable/gorm/gen_curdl.go.tpl
var CurdGen embed.FS

//go:embed enable/gorm/database.yaml
var Database embed.FS

//go:embed enable/gorm/database-alpha.yaml
var DatabaseAlpha embed.FS

//go:embed enable/gorm/database-release.yaml
var DatabaseRelease embed.FS

//go:embed docker/docker-compose-env.yml
var DockerComposeEnv embed.FS

//go:embed docker/docker/*
var DockerDir embed.FS

//go:embed .gitlab-ci.yml
var GitlabCI embed.FS

func CopyDirFromEmbedFs(embedFS embed.FS, src string, dest string) error {
	return fs.WalkDir(embedFS, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dest, relPath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, fs.ModePerm)
		} else {
			data, err := embedFS.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(targetPath, data, fs.ModePerm)
		}
	})
}

// 使用 embedFS 绑定文件，从 src 复制到 dest
func CopyFileFromEmbed(embedFS embed.FS, src, dest string) error {
	return fs.WalkDir(embedFS, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// 生成目标路径
		targetPath := filepath.Join(dest, relPath)
		filekit.EnsureDirExists(filepath.Dir(targetPath))

		// 如果是目录，创建目标目录
		if d.IsDir() {
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return err
			}
		} else {
			// 如果是文件，读取文件内容并写入目标文件
			data, err := embedFS.ReadFile(path)
			if err != nil {
				return err
			}
			if err := os.WriteFile(targetPath, data, os.ModePerm); err != nil {
				return err
			}
		}

		return nil
	})
}
