package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type fileSystem struct {
	path, index string
}

const separator = string(filepath.Separator)

func (d fileSystem) Open(name string) (http.File, error) {
	_, err := os.Stat(filepath.Join(d.path, strings.Split(name, separator)[1]))
	if os.IsNotExist(err) {
		name = separator + d.index
	}
	return http.Dir(d.path).Open(name)
}

func DirIndex(path, index string) http.FileSystem {
	return &fileSystem{
		path:  path,
		index: index,
	}
}

func Dir(path string) http.FileSystem {
	return DirIndex(path, "index.html")
}
