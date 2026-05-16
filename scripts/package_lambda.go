package main

import (
	"archive/zip"
	"flag"
	"io"
	"os"
	"path/filepath"
)

func main() {
	src := flag.String("src", "", "path to bootstrap binary")
	dst := flag.String("dst", "", "path to destination zip")
	flag.Parse()

	if *src == "" || *dst == "" {
		panic("src and dst are required")
	}

	if err := os.MkdirAll(filepath.Dir(*dst), 0o755); err != nil {
		panic(err)
	}

	zipFile, err := os.Create(*dst)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	srcFile, err := os.Open(*src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	w, err := zw.Create("bootstrap")
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(w, srcFile); err != nil {
		panic(err)
	}
}
