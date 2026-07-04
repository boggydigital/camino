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
	AbsDir int
	RelDir int
)

var (
	absPaths      map[AbsDir]string
	relNames      map[RelDir]string
	relAbsParents map[RelDir][]AbsDir
)

func Register(absolutePaths map[AbsDir]string, relativeNames map[RelDir]string, relativeAbsoluteParents map[RelDir][]AbsDir) error {

	absPaths = absolutePaths
	relNames = relativeNames
	relAbsParents = relativeAbsoluteParents

	for _, ap := range absolutePaths {
		if _, err := os.Stat(ap); os.IsNotExist(err) {
			if err = os.MkdirAll(ap, DefaultFileMode); err != nil {
				return err
			}
		}
	}

	for rp, relName := range relativeNames {

		if aps, ok := relativeAbsoluteParents[rp]; ok && len(aps) > 0 {

			for _, ap := range aps {
				if absPath, sure := absolutePaths[ap]; sure {

					absRelPath := filepath.Join(absPath, relName)

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

func GetAbs(absolutePath AbsDir) string {
	if ap, ok := absPaths[absolutePath]; ok {
		return ap
	}

	panic("abs path not registered")
}

func GetRel(relativePath RelDir, absolutePath AbsDir) string {

	if aps, ok := relAbsParents[relativePath]; ok {
		if !slices.Contains(aps, absolutePath) {
			panic("rel path not registered under abs path")
		}

		if relName, sure := relNames[relativePath]; sure {
			return filepath.Join(GetAbs(absolutePath), relName)
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

func ResolveAbsPaths(rootPath string, absDirNames map[AbsDir]string, overrides map[string]string) map[AbsDir]string {

	resolvedAbsPaths := make(map[AbsDir]string)

	for ad := range absDirNames {

		adn := absDirNames[ad]

		if ovd, ok := overrides[adn]; !ok {
			resolvedAbsPaths[ad] = filepath.Join(rootPath, adn)
		} else {
			resolvedAbsPaths[ad] = ovd
		}
	}

	return resolvedAbsPaths
}
