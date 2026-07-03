package camino

import "os"

const (
	permUrwGrwOr    os.FileMode = 0775
	DefaultFileMode             = permUrwGrwOr
)
