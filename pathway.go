package camino

type (
	AbsDir string
	RelDir string
)

type Pathway interface {
	AbsDirPath(absDir AbsDir) string
	AbsRelDirPath(relDir RelDir, absDir AbsDir) string
}
