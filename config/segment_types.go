package config

import (
	"errors"

	"github.com/LNKLEO/OMP/properties"
	"github.com/LNKLEO/OMP/runtime"
	"github.com/LNKLEO/OMP/segments"
)

// SegmentType the type of segment, for more information, see the constants
type SegmentType string

// SegmentWriter is the interface used to define what and if to write to the prompt
type SegmentWriter interface {
	Enabled() bool
	Template() string
	SetText(text string)
	Text() string
	Init(props properties.Properties, env runtime.Environment)
}

const (
	// Plain writes it without ornaments
	Plain SegmentStyle = "plain"
	// Powerline writes it Powerline style
	Powerline SegmentStyle = "powerline"
	// Accordion writes it Powerline style but collapses the segment when disabled instead of hiding
	Accordion SegmentStyle = "accordion"
	// Diamond writes the prompt shaped with a leading and trailing symbol
	Diamond SegmentStyle = "diamond"
	// ANGULAR writes which angular cli version us currently active
	ANGULAR SegmentType = "angular"
	// ARGOCD writes the current argocd context
	ARGOCD SegmentType = "argocd"
	// AURELIA writes which aurelia version is currently referenced in package.json
	AURELIA SegmentType = "aurelia"
	// AWS writes the active aws context
	AWS SegmentType = "aws"
	// AZ writes the Azure subscription info we're currently in
	AZ SegmentType = "az"
	// AZD writes the Azure Developer CLI environment info we're current in
	AZD SegmentType = "azd"
	// AZFUNC writes current AZ func version
	AZFUNC SegmentType = "azfunc"
	// BATTERY writes the battery percentage
	BATTERY SegmentType = "battery"
	// BAZEL writes the bazel version
	BAZEL SegmentType = "bazel"
	// Buf segment writes the active buf version
	BUF SegmentType = "buf"
	// CARBONINTENSITY writes the actual and forecast carbon intensity in gCO2/kWh
	CARBONINTENSITY SegmentType = "carbonintensity"
	// CMAKE writes the active cmake version
	CMAKE SegmentType = "cmake"
	// CMD writes the output of a shell command
	CMD SegmentType = "command"
	// CONNECTION writes a connection's information
	CONNECTION SegmentType = "connection"
	// DOCKER writes the docker context
	DOCKER SegmentType = "docker"
	// DOTNET writes which dotnet version is currently active
	DOTNET SegmentType = "dotnet"
	// EXECUTIONTIME writes the execution time of the last run command
	EXECUTIONTIME SegmentType = "executiontime"
	// EXIT writes the last exit code
	EXIT SegmentType = "exit"
	// GCP writes the active GCP context
	GCP SegmentType = "gcp"
	// GIT represents the git status and information
	GIT SegmentType = "git"
	// GITVERSION represents the gitversion information
	GITVERSION SegmentType = "gitversion"
	// GOLANG writes which go version is currently active
	GOLANG SegmentType = "go"
	// HASKELL segment
	HASKELL SegmentType = "haskell"
	// HELM segment
	HELM SegmentType = "helm"
	// JAVA writes the active java version
	JAVA SegmentType = "java"
	// KOTLIN writes the active kotlin version
	KOTLIN SegmentType = "kotlin"
	// KUBECTL writes the Kubernetes context we're currently in
	KUBECTL SegmentType = "kubectl"
	// LUA writes the active lua version
	LUA SegmentType = "lua"
	// MERCURIAL writes the Mercurial source control information
	MERCURIAL SegmentType = "mercurial"
	// MOJO writes the active version of Mojo and the name of the Magic virtual env
	MOJO SegmentType = "mojo"
	// MVN writes the active maven version
	MVN SegmentType = "mvn"
	// NODE writes which node version is currently active
	NODE SegmentType = "node"
	// npm version
	NPM SegmentType = "npm"
	// OS write os specific icon
	OS SegmentType = "os"
	// OWM writes the weather coming from openweatherdata
	OWM SegmentType = "owm"
	// PATH represents the current path segment
	PATH SegmentType = "path"
	// PROJECT
	PROJECT SegmentType = "project"
	// PYTHON writes the virtual env name
	PYTHON SegmentType = "python"
	// QUASAR writes the QUASAR version and context
	QUASAR SegmentType = "quasar"
	// REACT writes the current react version
	REACT SegmentType = "react"
	// ROOT writes root symbol
	ROOT SegmentType = "root"
	// RUST writes the cargo version information if cargo.toml is present
	RUST SegmentType = "rust"
	// SESSION represents the user info segment
	SESSION SegmentType = "session"
	// SHELL writes which shell we're currently in
	SHELL SegmentType = "shell"
	// STATUS writes the last know command status
	STATUS SegmentType = "status"
	// Subversion segment
	SVN SegmentType = "svn"
	// SYSTEMINFO writes system information (memory, cpu, load)
	SYSTEMINFO SegmentType = "sysinfo"
	// TALOSCTL writes the talosctl context
	TALOSCTL SegmentType = "talosctl"
	// TEXT writes a text
	TEXT SegmentType = "text"
	// TIME writes the current timestamp
	TIME SegmentType = "time"
	// WINREG queries the Windows registry.
	WINREG SegmentType = "winreg"
	// XMAKE write the xmake version if xmake.lua is present
	XMAKE SegmentType = "xmake"
)

// Segments contains all available prompt segment writers.
// Consumers of the library can also add their own segment writer.
var Segments = map[SegmentType]func() SegmentWriter{
	ANGULAR:         func() SegmentWriter { return &segments.Angular{} },
	ARGOCD:          func() SegmentWriter { return &segments.Argocd{} },
	AURELIA:         func() SegmentWriter { return &segments.Aurelia{} },
	AWS:             func() SegmentWriter { return &segments.Aws{} },
	AZ:              func() SegmentWriter { return &segments.Az{} },
	AZD:             func() SegmentWriter { return &segments.Azd{} },
	AZFUNC:          func() SegmentWriter { return &segments.AzFunc{} },
	BATTERY:         func() SegmentWriter { return &segments.Battery{} },
	BAZEL:           func() SegmentWriter { return &segments.Bazel{} },
	BUF:             func() SegmentWriter { return &segments.Buf{} },
	CARBONINTENSITY: func() SegmentWriter { return &segments.CarbonIntensity{} },
	CMAKE:           func() SegmentWriter { return &segments.Cmake{} },
	CMD:             func() SegmentWriter { return &segments.Cmd{} },
	CONNECTION:      func() SegmentWriter { return &segments.Connection{} },
	DOCKER:          func() SegmentWriter { return &segments.Docker{} },
	DOTNET:          func() SegmentWriter { return &segments.Dotnet{} },
	EXECUTIONTIME:   func() SegmentWriter { return &segments.Executiontime{} },
	EXIT:            func() SegmentWriter { return &segments.Status{} },
	GCP:             func() SegmentWriter { return &segments.Gcp{} },
	GIT:             func() SegmentWriter { return &segments.Git{} },
	GITVERSION:      func() SegmentWriter { return &segments.GitVersion{} },
	GOLANG:          func() SegmentWriter { return &segments.Golang{} },
	HASKELL:         func() SegmentWriter { return &segments.Haskell{} },
	HELM:            func() SegmentWriter { return &segments.Helm{} },
	JAVA:            func() SegmentWriter { return &segments.Java{} },
	KOTLIN:          func() SegmentWriter { return &segments.Kotlin{} },
	KUBECTL:         func() SegmentWriter { return &segments.Kubectl{} },
	LUA:             func() SegmentWriter { return &segments.Lua{} },
	MERCURIAL:       func() SegmentWriter { return &segments.Mercurial{} },
	MOJO:            func() SegmentWriter { return &segments.Mojo{} },
	MVN:             func() SegmentWriter { return &segments.Mvn{} },
	NODE:            func() SegmentWriter { return &segments.Node{} },
	NPM:             func() SegmentWriter { return &segments.Npm{} },
	OS:              func() SegmentWriter { return &segments.Os{} },
	OWM:             func() SegmentWriter { return &segments.Owm{} },
	PATH:            func() SegmentWriter { return &segments.Path{} },
	PROJECT:         func() SegmentWriter { return &segments.Project{} },
	PYTHON:          func() SegmentWriter { return &segments.Python{} },
	QUASAR:          func() SegmentWriter { return &segments.Quasar{} },
	REACT:           func() SegmentWriter { return &segments.React{} },
	ROOT:            func() SegmentWriter { return &segments.Root{} },
	RUST:            func() SegmentWriter { return &segments.Rust{} },
	SESSION:         func() SegmentWriter { return &segments.Session{} },
	SHELL:           func() SegmentWriter { return &segments.Shell{} },
	STATUS:          func() SegmentWriter { return &segments.Status{} },
	SVN:             func() SegmentWriter { return &segments.Svn{} },
	SYSTEMINFO:      func() SegmentWriter { return &segments.SystemInfo{} },
	TALOSCTL:        func() SegmentWriter { return &segments.TalosCTL{} },
	TEXT:            func() SegmentWriter { return &segments.Text{} },
	TIME:            func() SegmentWriter { return &segments.Time{} },
	WINREG:          func() SegmentWriter { return &segments.WindowsRegistry{} },
	XMAKE:           func() SegmentWriter { return &segments.XMake{} },
}

func (segment *Segment) MapSegmentWithWriter(env runtime.Environment) error {
	segment.env = env

	if segment.Properties == nil {
		segment.Properties = make(properties.Map)
	}

	f, ok := Segments[segment.Type]
	if !ok {
		return errors.New("unable to map writer")
	}

	writer := f()
	wrapper := &properties.Wrapper{
		Properties: segment.Properties,
	}

	writer.Init(wrapper, env)
	segment.writer = writer

	return nil
}
