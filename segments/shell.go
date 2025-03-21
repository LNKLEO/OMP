package segments

import (
	"strings"

	"github.com/LNKLEO/OMP/properties"
)

type Shell struct {
	base

	Name    string
	Version string
}

const (
	// MappedShellNames allows for custom text in place of shell names
	MappedShellNames properties.Property = "mapped_shell_names"
)

func (s *Shell) Template() string {
	return NameTemplate
}

func (s *Shell) Enabled() bool {
	mappedNames := s.props.GetKeyValueMap(MappedShellNames, make(map[string]string))
	s.Name = s.env.Shell()
	s.Version = s.env.Flags().ShellVersion
	for key, val := range mappedNames {
		if strings.EqualFold(s.Name, key) {
			s.Name = val
			break
		}
	}
	return true
}
