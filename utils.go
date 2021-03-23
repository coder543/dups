package dups

import (
	"fmt"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

// CleanPath replaces \ with / in a path
func CleanPath(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

// createBar creates a new progress bar with a custom template
func createBar(limit int64, fullHash bool) *pb.ProgressBar {
	prefix := "initial pass"
	if fullHash {
		prefix = "verification"
	}

	tmpl := fmt.Sprintf(
		`%s {{ blue "Progress:" }} {{ bar . "[" "=" (cycle . ">") "." "]"}} {{speed . | green }} {{ rtime . | green }} {{percent . | green}}`,
		prefix,
	)

	bar := pb.ProgressBarTemplate(tmpl).Start64(limit).Set(pb.Bytes, true)
	return bar
}
