package prompt

import (
	"github.com/LNKLEO/OMP/config"
)

func (e *Engine) RPrompt() string {
	var rprompt *config.Block

	for _, block := range e.Config.Blocks {
		if block.Type != config.RPrompt {
			continue
		}

		rprompt = block
		break
	}

	if rprompt == nil {
		return ""
	}

	text, length := e.writeBlockSegments(rprompt)

	// do not print anything when we don't have any text
	if length == 0 {
		return ""
	}

	e.rpromptLength = length

	return text
}
