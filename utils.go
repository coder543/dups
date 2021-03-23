package dups

import (
	"strings"

	"github.com/cheggaaa/pb/v3"
)

// CleanPath replaces \ with / in a path
func CleanPath(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

// createBar creates a new progress bar with a custom template
func createBar(limit int64) *pb.ProgressBar {
	tmpl := `{{ blue "Progress:" }} {{ bar . "[" "=" (cycle . ">") "." "]"}} {{speed . | green }} {{ rtime . | green }} {{percent . | green}}`
	bar := pb.ProgressBarTemplate(tmpl).Start64(limit).Set(pb.Bytes, true)
	return bar
}
