package segments

import (
	"github.com/LNKLEO/OMP/platform"
	"github.com/LNKLEO/OMP/properties"
)

type Elixir struct {
	language
}

func (e *Elixir) Template() string {
	return languageTemplate
}

func (e *Elixir) Init(props properties.Properties, env platform.Environment) {
	e.language = language{
		env:        env,
		props:      props,
		extensions: []string{"*.ex", "*.exs"},
		commands: []*cmd{
			{
				executable: "asdf",
				args:       []string{"current", "elixir"},
				regex:      `elixir\s+(?P<version>((?P<major>[0-9]+).(?P<minor>[0-9]+).(?P<patch>[0-9]+)))[^\s]*\s+`,
			},
			{
				executable: "elixir",
				args:       []string{"--version"},
				regex:      `Elixir (?P<version>((?P<major>[0-9]+).(?P<minor>[0-9]+).(?P<patch>[0-9]+)))`,
			},
		},
		versionURLTemplate: "https://github.com/elixir-lang/elixir/releases/tag/v{{ .Full }}",
	}
}

func (e *Elixir) Enabled() bool {
	return e.language.Enabled()
}
