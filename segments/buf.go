package segments

import (
	"github.com/LNKLEO/OMP/platform"
	"github.com/LNKLEO/OMP/properties"
)

type Buf struct {
	language
}

func (b *Buf) Template() string {
	return languageTemplate
}

func (b *Buf) Init(props properties.Properties, env platform.Environment) {
	b.language = language{
		env:        env,
		props:      props,
		extensions: []string{"buf.yaml", "buf.gen.yaml", "buf.work.yaml"},
		commands: []*cmd{
			{
				executable: "buf",
				args:       []string{"--version"},
				regex:      `(?:(?P<version>((?P<major>[0-9]+).(?P<minor>[0-9]+).(?P<patch>[0-9]+))))`,
			},
		},
		versionURLTemplate: "https://github.com/bufbuild/buf/releases/tag/v{{.Full}}",
	}
}

func (b *Buf) Enabled() bool {
	return b.language.Enabled()
}
