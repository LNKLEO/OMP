package shell

import (
	_ "embed"
)

//go:embed scripts/omp.zsh
var zshInit string

func (f Feature) Zsh() Code {
	switch f {
	case CursorPositioning:
		return unixCursorPositioning
	case Tooltips:
		return "enable_omptooltips"
	case Transient:
		return "_omp_create_widget zle-line-init _omp_zle-line-init"
	case FTCSMarks:
		return unixFTCSMarks
	case RPrompt, OMPGit, Azure, LineError, Jobs:
		fallthrough
	default:
		return ""
	}
}
