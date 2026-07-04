package camino

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	tarGzExt = ".tar.gz"
)

const Layout = "2006-01-02-15-04-05"

func TimedTarGzFilename() string {
	return time.Now().Format(Layout) + tarGzExt
}

func Compress(src, dst string) error {

	exportedPath := filepath.Join(dst, TimedTarGzFilename())

	if _, err := os.Stat(exportedPath); os.IsExist(err) {
		return err
	}

	file, err := os.Create(exportedPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	tarWalker := func(path string, fi os.FileInfo, err error) error {

		if fi.Mode().IsDir() {
			return nil
		}

		// this takes care of linked files that are problematic for tar
		if !fi.Mode().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if len(relPath) == 0 {
			return nil
		}

		rcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer rcFile.Close()

		var h *tar.Header
		if h, err = tar.FileInfoHeader(fi, relPath); err != nil {
			return err
		}

		h.Name = relPath
		if err = tw.WriteHeader(h); err != nil {
			return err
		}

		if _, err = io.Copy(tw, rcFile); err != nil {
			return err
		}
		return nil
	}

	if err = filepath.Walk(src, tarWalker); err != nil {
		return err
	}

	return nil

}
