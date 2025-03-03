package shell

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed scripts/omp.bash
var bashInit string

func (f Feature) Bash() Code {
	switch f {
	case CursorPositioning:
		return unixCursorPositioning
	case FTCSMarks:
		return unixFTCSMarks
	case PromptMark, RPrompt, Git, Azure, LineError, Jobs, Tooltips, Transient:
		fallthrough
	default:
		return ""
	}
}

func QuotePosixStr(str string) string {
	if len(str) == 0 {
		return "''"
	}

	return fmt.Sprintf("$'%s'", strings.NewReplacer(`\`, `\\`, "'", `\'`).Replace(str))
}
