package camino

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	daysToPreserveFiles = 30
	numFilesToPreserve  = 15
)

func CleanupTimed(dir AbsDir, delete bool) error {

	absDirPath := GetAbs(dir)

	d, err := os.Open(absDirPath)
	if err != nil {
		return err
	}
	defer d.Close()

	filenames, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	earliest := time.Now().Add(-daysToPreserveFiles * 24 * time.Hour)
	oldFiles := make([]string, 0)
	currentFiles := make([]string, 0)

	for _, fn := range filenames {

		fnse := fn
		for filepath.Ext(fnse) != "" {
			fnse = strings.TrimSuffix(fnse, filepath.Ext(fnse))
		}
		ft, perr := time.Parse(Layout, fnse)
		if perr != nil {
			continue
		}

		if ft.After(earliest) {
			currentFiles = append(currentFiles, fn)
			continue
		}

		oldFiles = append(oldFiles, fn)
	}

	// first, delete old backups
	if len(oldFiles) > 0 && delete {

		// never delete all backups, leave the latest file as the current backup
		if len(oldFiles) == len(filenames) {
			if err = os.Rename(oldFiles[len(oldFiles)-1], TimestampedTarGzFilename()); err != nil {
				return err
			}
			oldFiles = oldFiles[:len(oldFiles)-1]
		}

		for _, fn := range oldFiles {
			filename := filepath.Join(absDirPath, fn)
			if err = os.Remove(filename); err != nil {
				return err
			}
		}
	}

	// second, trim backups to the specified count
	if len(currentFiles) > numFilesToPreserve {

		slices.Sort(currentFiles)

		for ii := 0; ii < (len(currentFiles) - numFilesToPreserve); ii++ {
			filename := filepath.Join(absDirPath, currentFiles[ii])
			if err = os.Remove(filename); err != nil {
				return err
			}
		}

	}

	return nil
}
