package shell

type Feature byte

const (
	Jobs Feature = iota
	Azure
	Git
	LineError
	Tooltips
	Transient
	FTCSMarks
	PromptMark
	RPrompt
	CursorPositioning
)

type Features []Feature

func (f Features) Lines(shell string) Lines {
	var lines Lines

	for _, feature := range f {
		var code Code

		switch shell {
		case PWSH, PWSH5:
			code = feature.Pwsh()
		case ZSH:
			code = feature.Zsh()
		case BASH:
			code = feature.Bash()
		case CMD:
			code = feature.Cmd()
		}

		if len(code) > 0 {
			lines = append(lines, code)
		}
	}

	return lines
}
