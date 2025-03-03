package shell

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed scripts/omp.ps1
var pwshInit string

func (f Feature) Pwsh() Code {
	switch f {
	case Tooltips:
		return "Enable-OMPTooltips"
	case LineError:
		return "Enable-OMPLineError"
	case Transient:
		return "Enable-OMPTransientPrompt"
	case Jobs:
		return "$global:_ompJobCount = $true"
	case Azure:
		return "$global:_ompAzure = $true"
	case Git:
		return "$global:_ompGit = $true"
	case FTCSMarks:
		return "$global:_ompFTCSMarks = $true"
	case PromptMark, RPrompt, CursorPositioning:
		fallthrough
	default:
		return ""
	}
}

func quotePwshStr(str string) string {
	if len(str) == 0 {
		return "''"
	}

	return fmt.Sprintf("'%s'", strings.ReplaceAll(str, "'", "''"))
}
