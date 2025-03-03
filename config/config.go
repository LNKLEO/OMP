package config

import (
	"github.com/LNKLEO/OMP/color"
	"github.com/LNKLEO/OMP/runtime"
	"github.com/LNKLEO/OMP/segments"
	"github.com/LNKLEO/OMP/shell"
	"github.com/LNKLEO/OMP/template"
)

const (
	JSON string = "json"
	YAML string = "yaml"
	TOML string = "toml"

	Version = 3
)

// Config holds all the theme for rendering the prompt
type Config struct {
	Palette                 color.Palette   `json:"palette,omitempty" toml:"palette,omitempty"`
	DebugPrompt             *Segment        `json:"debug_prompt,omitempty" toml:"debug_prompt,omitempty"`
	Var                     map[string]any  `json:"var,omitempty" toml:"var,omitempty"`
	Palettes                *color.Palettes `json:"palettes,omitempty" toml:"palettes,omitempty"`
	ValidLine               *Segment        `json:"valid_line,omitempty" toml:"valid_line,omitempty"`
	SecondaryPrompt         *Segment        `json:"secondary_prompt,omitempty" toml:"secondary_prompt,omitempty"`
	TransientPrompt         *Segment        `json:"transient_prompt,omitempty" toml:"transient_prompt,omitempty"`
	ErrorLine               *Segment        `json:"error_line,omitempty" toml:"error_line,omitempty"`
	TerminalBackground      color.Ansi      `json:"terminal_background,omitempty" toml:"terminal_background,omitempty"`
	origin                  string
	PWD                     string                 `json:"pwd,omitempty" toml:"pwd,omitempty"`
	AccentColor             color.Ansi             `json:"accent_color,omitempty" toml:"accent_color,omitempty"`
	Output                  string                 `json:"-" toml:"-"`
	ConsoleTitleTemplate    string                 `json:"console_title_template,omitempty" toml:"console_title_template,omitempty"`
	Format                  string                 `json:"-" toml:"-"`
	Cycle                   color.Cycle            `json:"cycle,omitempty" toml:"cycle,omitempty"`
	Blocks                  []*Block               `json:"blocks,omitempty" toml:"blocks,omitempty"`
	Tooltips                []*Segment             `json:"tooltips,omitempty" toml:"tooltips,omitempty"`
	Version                 int                    `json:"version" toml:"version"`
	ShellIntegration        bool                   `json:"shell_integration,omitempty" toml:"shell_integration,omitempty"`
	MigrateGlyphs           bool                   `json:"-" toml:"-"`
	PatchPwshBleed          bool                   `json:"patch_pwsh_bleed,omitempty" toml:"patch_pwsh_bleed,omitempty"`
	EnableCursorPositioning bool                   `json:"enable_cursor_positioning,omitempty" toml:"enable_cursor_positioning,omitempty"`
	FinalSpace              bool                   `json:"final_space,omitempty" toml:"final_space,omitempty"`
}

func (cfg *Config) MakeColors(env runtime.Environment) color.String {
	cacheDisabled := env.Getenv("OMP_CACHE_DISABLED") == "1"
	return color.MakeColors(cfg.getPalette(), !cacheDisabled, cfg.AccentColor, env)
}

func (cfg *Config) getPalette() color.Palette {
	if cfg.Palettes == nil {
		return cfg.Palette
	}

	tmpl := &template.Text{
		Template: cfg.Palettes.Template,
	}

	key, err := tmpl.Render()
	if err != nil {
		return cfg.Palette
	}

	palette, ok := cfg.Palettes.List[key]
	if !ok {
		return cfg.Palette
	}

	for key, color := range cfg.Palette {
		if _, ok := palette[key]; ok {
			continue
		}

		palette[key] = color
	}

	return palette
}

func (cfg *Config) Features(env runtime.Environment) shell.Features {
	var feats shell.Features

	if cfg.TransientPrompt != nil {
		feats = append(feats, shell.Transient)
	}

	if cfg.ShellIntegration {
		feats = append(feats, shell.FTCSMarks)
	}

	if cfg.ErrorLine != nil || cfg.ValidLine != nil {
		feats = append(feats, shell.LineError)
	}

	if len(cfg.Tooltips) > 0 {
		feats = append(feats, shell.Tooltips)
	}

	for i, block := range cfg.Blocks {
		if (i == 0 && block.Newline) && cfg.EnableCursorPositioning {
			feats = append(feats, shell.CursorPositioning)
		}

		if block.Type == RPrompt {
			feats = append(feats, shell.RPrompt)
		}

		for _, segment := range block.Segments {
			if segment.Type == AZ {
				source := segment.Properties.GetString(segments.Source, segments.FirstMatch)
				if source == segments.Pwsh || source == segments.FirstMatch {
					feats = append(feats, shell.Azure)
				}
			}

			if segment.Type == GIT {
				source := segment.Properties.GetString(segments.Source, segments.Cli)
				if source == segments.Pwsh {
					feats = append(feats, shell.Git)
				}
			}
		}
	}

	return feats
}
