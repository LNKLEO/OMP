package segments

import (
	"github.com/LNKLEO/OMP/platform"
	"github.com/LNKLEO/OMP/properties"
)

type Vala struct {
	language
}

func (v *Vala) Template() string {
	return languageTemplate
}

func (v *Vala) Init(props properties.Properties, env platform.Environment) {
	v.language = language{
		env:        env,
		props:      props,
		extensions: []string{"*.vala"},
		commands: []*cmd{
			{
				executable: "vala",
				args:       []string{"--version"},
				regex:      `Vala (?P<version>((?P<major>[0-9]+).(?P<minor>[0-9]+).(?P<patch>[0-9]+)))`,
			},
		},
		versionURLTemplate: "https://gitlab.gnome.org/GNOME/vala/raw/{{ .Major }}.{{ .Minor }}/NEWS",
	}
}

func (v *Vala) Enabled() bool {
	return v.language.Enabled()
}
