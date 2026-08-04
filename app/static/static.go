// Package static contains code for serving static file assets
package static

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed *
var Files embed.FS

type noDirFS struct {
	fs http.FileSystem
}

func (nfs noDirFS) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, ".go") {
		return nil, fs.ErrNotExist
	}
	f, err := nfs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if stat.IsDir() {
		_ = f.Close()
		return nil, fs.ErrNotExist
	}
	return f, nil
}

func Handler() http.Handler {
	return http.FileServer(noDirFS{http.FS(Files)})
}

func Exists(name string) bool {
	name = strings.TrimPrefix(name, "/")
	if name == "" || strings.HasSuffix(name, ".go") {
		return false
	}
	f, err := Files.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil || stat.IsDir() {
		return false
	}
	return true
}
