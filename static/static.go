package static

import (
	"embed"
	"io/fs"
)

//go:embed assets
var Assets embed.FS

// StaticFS 返回静态资源的子文件系统,用于HTTP静态文件服务
func StaticFS() (fs.FS, error) {
	return fs.Sub(Assets, "assets/www")
}
