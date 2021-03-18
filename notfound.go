package main

import (
	"embed"
	"errors"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
)

var ErrDir = errors.New("path is dir")

func tryRead(fs embed.FS, prefix, requestedPath string, w http.ResponseWriter) error {
	f, err := fs.Open(path.Join(prefix, requestedPath))
	if err != nil {
		return err
	}
	defer f.Close()

	// Goのfs.Openはディレクトリを読みこもとうしてもエラーにはならないがここでは邪魔なのでエラー扱いにする
	stat, _ := f.Stat()
	if stat.IsDir() {
		return ErrDir
	}

	contentType := mime.TypeByExtension(filepath.Ext(requestedPath))
	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, f)
	return err
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// まずはリクエストされた通りにファイルを探索
	err := tryRead(assets, "frontend/dist", r.URL.Path, w)
	if err == nil {
		return
	}
	// 見つからなければindex.htmlを返す
	err = tryRead(assets, "frontend/dist", "index.html", w)
	if err != nil {
		panic(err)
	}
}
