package segments

import (
	"github.com/LNKLEO/OMP/platform"
	"github.com/LNKLEO/OMP/properties"
	"github.com/LNKLEO/OMP/regex"
)

type Session struct {
	props properties.Properties
	env   platform.Environment
	// text  string

	SSHSession bool

	// Deprecated
	DefaultUserName string
}

func (s *Session) Enabled() bool {
	s.SSHSession = s.activeSSHSession()
	return true
}

func (s *Session) Template() string {
	return " {{ if .SSHSession }}\ueba9 {{ end }}{{ .UserName }}@{{ .HostName }} "
}

func (s *Session) Init(props properties.Properties, env platform.Environment) {
	s.props = props
	s.env = env
}

func (s *Session) activeSSHSession() bool {
	keys := []string{
		"SSH_CONNECTION",
		"SSH_CLIENT",
	}

	for _, key := range keys {
		content := s.env.Getenv(key)
		if content != "" {
			return true
		}
	}

	if s.env.Platform() == platform.WINDOWS {
		return false
	}

	whoAmI, err := s.env.RunCommand("who", "am", "i")
	if err != nil {
		return false
	}

	return regex.MatchString(`\(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\)`, whoAmI)
}
