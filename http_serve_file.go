package camino

import (
	"net/http"
	"os"
	"slices"
	"strconv"
)

type ServeOption int

const (
	NoCacheControl ServeOption = iota
	NoContentDisposition
	NoContentLength
	NoBinaryContentType
	NoLastModified
)

func ServeFile(absPath string, w http.ResponseWriter, r *http.Request, opts ...ServeOption) {

	if fi, err := os.Stat(absPath); err == nil {

		if !slices.Contains(opts, NoCacheControl) {
			w.Header().Set("Cache-Control", "max-age=31536000")
		}
		if !slices.Contains(opts, NoContentDisposition) {
			w.Header().Set("Content-Disposition", "attachment; filename=\""+fi.Name()+"\"")
		}
		if !slices.Contains(opts, NoContentLength) {
			w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
		}
		if !slices.Contains(opts, NoBinaryContentType) {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
		if !slices.Contains(opts, NoLastModified) {
			w.Header().Set("Last-Modified", fi.ModTime().Format(http.TimeFormat))
		}

		http.ServeFile(w, r, absPath)
	} else {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
