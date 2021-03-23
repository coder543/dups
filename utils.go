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

const (
	BAR_INITIAL = iota
	BAR_FULLHASH
	BAR_LINKING
	BAR_DELETING
)

// createBar creates a new progress bar with a custom template
func createBar(limit int64, barType int) *pb.ProgressBar {
	var prefix string
	var bytes bool
	switch barType {
	case BAR_INITIAL:
		prefix = "Initial Pass"
		bytes = true
	case BAR_FULLHASH:
		prefix = "Verification"
		bytes = true
	case BAR_LINKING:
		prefix = "Linking"
	case BAR_DELETING:
		prefix = "Deleting"
	}

	tmpl := fmt.Sprintf(
		`{{ blue "%s:" }} {{ bar . "[" "=" (cycle . ">") "." "]"}} {{speed . | green }} {{ rtime . | blue }} {{percent . | green}}`,
		prefix,
	)

	bar := pb.ProgressBarTemplate(tmpl).Start64(limit)

	if bytes {
		bar = bar.Set(pb.Bytes, true)
	}

	return bar
}
