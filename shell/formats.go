package shell

type Formats struct {
	Escape     string
	Left       string
	Linechange string
	ClearBelow string
	ClearLine  string

	Title string

	SaveCursorPosition    string
	RestoreCursorPosition string

	Osc99 string
	Osc7  string
	Osc51 string

	EscapeSequences map[rune]string

	HyperlinkStart  string
	HyperlinkCenter string
	HyperlinkEnd    string
}

func GetFormats(shell string) *Formats {
	var formats *Formats

	switch shell {
	case BASH:
		formats = &Formats{
			Escape:                "\\[%s\\]",
			Linechange:            "\\[\x1b[%d%s\\]",
			Left:                  "\\[\x1b[%dD\\]",
			ClearBelow:            "\\[\x1b[0J\\]",
			ClearLine:             "\\[\x1b[K\\]",
			SaveCursorPosition:    "\\[\x1b7\\]",
			RestoreCursorPosition: "\\[\x1b8\\]",
			Title:                 "\\[\x1b]0;%s\007\\]",
			HyperlinkStart:        "\\[\x1b]8;;",
			HyperlinkCenter:       "\x1b\\\\\\]",
			HyperlinkEnd:          "\\[\x1b]8;;\x1b\\\\\\]",
			Osc99:                 "\\[\x1b]9;9;%s\x1b\\\\\\]",
			Osc7:                  "\\[\x1b]7;file://%s/%s\x1b\\\\\\]",
			Osc51:                 "\\[\x1b]51;A;%s@%s:%s\x1b\\\\\\]",
			EscapeSequences: map[rune]string{
				'\\': `\\`,
			},
		}
	case ZSH:
		formats = &Formats{
			Escape:                "%%{%s%%}",
			Linechange:            "%%{\x1b[%d%s%%}",
			Left:                  "%%{\x1b[%dD%%}",
			ClearBelow:            "%{\x1b[0J%}",
			ClearLine:             "%{\x1b[K%}",
			SaveCursorPosition:    "%{\x1b7%}",
			RestoreCursorPosition: "%{\x1b8%}",
			Title:                 "%%{\x1b]0;%s\007%%}",
			HyperlinkStart:        "%{\x1b]8;;",
			HyperlinkCenter:       "\x1b\\%}",
			HyperlinkEnd:          "%{\x1b]8;;\x1b\\%}",
			Osc99:                 "%%{\x1b]9;9;%s\x1b\\%%}",
			Osc7:                  "%%{\x1b]7;file://%s/%s\x1b\\%%}",
			Osc51:                 "%%{\x1b]51;A%s@%s:%s\x1b\\%%}",
			EscapeSequences: map[rune]string{
				'%': "%%",
			},
		}
	default:
		formats = &Formats{
			Escape:                "%s",
			Linechange:            "\x1b[%d%s",
			Left:                  "\x1b[%dD",
			ClearBelow:            "\x1b[0J",
			ClearLine:             "\x1b[K",
			SaveCursorPosition:    "\x1b7",
			RestoreCursorPosition: "\x1b8",
			Title:                 "\x1b]0;%s\007",
			HyperlinkStart:  "\x1b]8;;",
			HyperlinkCenter: "\x1b\\",
			HyperlinkEnd:    "\x1b]8;;\x1b\\",
			Osc99:           "\x1b]9;9;%s\x1b\\",
			Osc7:            "\x1b]7;file://%s/%s\x1b\\",
			Osc51:           "\x1b]51;A%s@%s:%s\x1b\\",
		}
	}

	return formats
}
