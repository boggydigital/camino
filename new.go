package camino

import "path/filepath"

type (
	AbsPath int
	RelPath int
)

var (
	relParentPath map[RelPath]AbsPath
	relPaths      map[RelPath]string
	absPaths      map[AbsPath]string
)

func Register(ap AbsPath, absPath string, rps map[RelPath]string) {
	if absPaths == nil {
		absPaths = make(map[AbsPath]string)
	}

	absPaths[ap] = absPath

	if len(rps) == 0 {
		return
	}

	if relParentPath == nil {
		relParentPath = make(map[RelPath]AbsPath)
	}

	if relPaths == nil {
		relPaths = make(map[RelPath]string)
	}

	for rp, relPath := range relPaths {
		relPaths[rp] = relPath
		relParentPath[rp] = ap
	}
}

func GetAbs(ap AbsPath) string {
	if absPath, ok := absPaths[ap]; ok {
		return absPath
	}

	panic("abs path not registered")
}

func GetRel(rp RelPath) string {
	if ap, ok := relParentPath[rp]; ok {
		if rpp, sure := relPaths[rp]; sure {
			return filepath.Join(GetAbs(ap), rpp)
		}
	}

	panic("rel path not registered")
}
