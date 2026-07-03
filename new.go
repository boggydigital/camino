package camino

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type (
	AbsPath int
	RelPath int
)

var (
	absPaths    map[AbsPath]string
	relPaths    map[RelPath]string
	relAbsPaths map[RelPath][]AbsPath
)

func Register(absolutePaths map[AbsPath]string, relativePaths map[RelPath]string, relativeAbsolutePaths map[RelPath][]AbsPath) error {

	absPaths = absolutePaths
	relPaths = relativePaths
	relAbsPaths = relativeAbsolutePaths

	for _, ap := range absolutePaths {
		if _, err := os.Stat(ap); os.IsNotExist(err) {
			if err = os.MkdirAll(ap, DefaultFileMode); err != nil {
				return err
			}
		}
	}

	for rp, relPath := range relativePaths {

		if aps, ok := relativeAbsolutePaths[rp]; ok && len(aps) > 0 {

			for _, ap := range aps {
				if absPath, sure := absolutePaths[ap]; sure {

					absRelPath := filepath.Join(absPath, relPath)

					if _, err := os.Stat(absRelPath); os.IsNotExist(err) {
						if err = os.MkdirAll(absRelPath, DefaultFileMode); err != nil {
							return err
						}
					}

				} else {
					return errors.New("unknown absolute path in relative absolute paths")
				}
			}

		} else {
			return errors.New("relative absolute paths not set")
		}

	}

	return nil
}

func GetAbs(absolutePath AbsPath) string {
	if ap, ok := absPaths[absolutePath]; ok {
		return ap
	}

	panic("abs path not registered")
}

func GetRel(relativePath RelPath, absolutePath AbsPath) string {

	if aps, ok := relAbsPaths[relativePath]; ok {
		if !slices.Contains(aps, absolutePath) {
			panic("rel path not registered under abs path")
		}

		if relPath, sure := relPaths[relativePath]; sure {
			return filepath.Join(GetAbs(absolutePath), relPath)
		}

		panic("rel path not registered")
	}

	panic("abs paths not registered for rel path")
}

func ReadOverrides(path string) (map[string]string, error) {

	directoriesFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer directoriesFile.Close()

	directories := make(map[string]string)

	scanner := bufio.NewScanner(directoriesFile)
	for scanner.Scan() {

		line := scanner.Text()

		if absDir, absDirPath, ok := strings.Cut(line, "="); ok {
			directories[absDir] = absDirPath
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return directories, nil
}
