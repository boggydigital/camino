package camino

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type setPathResolver struct {
	absDirPaths map[AbsDir]string
}

func (spr *setPathResolver) AbsDirPath(absDir AbsDir) string {
	if adp, ok := spr.absDirPaths[absDir]; ok {
		return adp
	} else {
		panic("pathways: unset abs dir " + string(absDir))
	}
}

func (spr *setPathResolver) AbsRelDirPath(relDir RelDir, absDir AbsDir) string {
	return filepath.Join(spr.AbsDirPath(absDir), string(relDir))
}

func NewSet(absDirPaths map[AbsDir]string) (Pathway, error) {

	for _, absDirPath := range absDirPaths {
		if _, err := os.Stat(absDirPath); err != nil {
			return nil, err
		}
	}

	spr := &setPathResolver{absDirPaths: absDirPaths}

	return spr, nil
}

func ReadSet(path string) (Pathway, error) {

	setFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer setFile.Close()

	absDirPaths := make(map[AbsDir]string)

	scanner := bufio.NewScanner(setFile)
	for scanner.Scan() {

		line := scanner.Text()

		if absDir, absDirPath, ok := strings.Cut(line, "="); ok {
			absDirPaths[AbsDir(absDir)] = absDirPath
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return NewSet(absDirPaths)
}
