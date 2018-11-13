package funcs

import (
	"testing"
)

func TestCacheReportFiles(t *testing.T) {
	InitChecksum()
	files := map[string]int{
		"/etc/passwd":   0,
		"/etc/my.cnf":   0,
		"/etc/my.cnf.d": 1,
	}
	SetReportFiles(files)
}
