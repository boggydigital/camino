package camino

import (
	"os"
	"path/filepath"
)

type rootPathResolver struct {
	rooDir string
}

func (rpr *rootPathResolver) AbsDirPath(absDir AbsDir) string {
	return filepath.Join(rpr.rooDir, string(absDir))
}

func (rpr *rootPathResolver) AbsRelDirPath(relDir RelDir, absDir AbsDir) string {
	return filepath.Join(rpr.AbsDirPath(absDir), string(relDir))
}

func NewRoot(rootDir string) (Pathway, error) {

	if _, err := os.Stat(rootDir); err != nil {
		return nil, err
	}

	rpr := &rootPathResolver{rooDir: rootDir}

	return rpr, nil
}
