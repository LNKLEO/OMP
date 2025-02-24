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
	case OMPGit:
		return "$global:_ompOMPGit = $true"
	case FTCSMarks:
		return "$global:_ompFTCSMarks = $true"
	case RPrompt, CursorPositioning:
		fallthrough
	default:
		return ""
	}
}

func quotePwshOrElvishStr(str string) string {
	if len(str) == 0 {
		return "''"
	}

	return fmt.Sprintf("'%s'", strings.ReplaceAll(str, "'", "''"))
}
