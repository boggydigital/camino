package camino

import (
	"slices"
	"unicode"
)

// https://en.wikipedia.org/wiki/Filename#Problematic_characters
var problematicChars = []rune{
	'/',
	'\\',
	'?',
	'%',
	'*',
	':',
	'|',
	'"',
	'<',
	'>',
	'.',
	',',
	';',
	'=',
	' '}

// https://en.wikipedia.org/wiki/Control_character#In_ASCII
var asciiControlChars = []rune{
	'\x00', '\x01', '\x02', '\x03', '\x04',
	'\x05', '\x06', '\x07', '\x08', '\x09',
	'\x0a', '\x0b', '\x0c', '\x0d', '\x0e', '\x0f',
	'\x10', '\x11', '\x12', '\x13', '\x14',
	'\x15', '\x16', '\x17', '\x18', '\x19',
	'\x1a', '\x1b', '\x1c', '\x1d', '\x1e', '\x1f',
	'\x7f'}

// Go is not supported on DOS: https://go.dev/dl/, no need to filter the set below:
//
// In addition, in Windows and DOS utilities, some words are also reserved
// and cannot be used as filenames.[19] For example, DOS device files:[21]
// CON, CONIN$, CONOUT$, PRN, AUX, CLOCK$, NUL
// COM0, COM1, COM2, COM3, COM4, COM5, COM6, COM7, COM8, COM9[7]
// LPT0, LPT1, LPT2, LPT3, LPT4, LPT5, LPT6, LPT7, LPT8, LPT9[7]
// LST (only in 86-DOS and DOS 1.xx)
// KEYBD$, SCREEN$ (only in multitasking MS-DOS 4.0)
// $IDLE$ (only in Concurrent DOS 386, Multiuser DOS and DR DOS 5.0 and higher)
// CONFIG$ (only in MS-DOS 7.0-8.0)

// NTFS filenames that are used internally include:
var ntfsIntFilenames = []string{
	"$Mft",
	"$MftMirr",
	"$LogFile",
	"$Volume",
	"$AttrDef",
	"$Bitmap",
	"$Boot",
	"$BadClus",
	"$Secure",
	"$Upcase",
	"$Extend",
	"$Quota",
	"$ObjId",
	"$Reparse",
}

const (
	safeStr = "_"
)

func Sanitize(s string) string {
	sans := ""
	for _, ch := range s {
		if slices.Contains(problematicChars, ch) {
			sans += safeStr
		} else if slices.Contains(asciiControlChars, ch) {
			sans += safeStr
		} else if unicode.IsSpace(ch) {
			sans += safeStr
		} else {
			sans += string(ch)
		}
	}

	if slices.Contains(ntfsIntFilenames, sans) {
		sans = safeStr + sans
	}

	return sans
}
