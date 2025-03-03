package segments

import (
	"fmt"
	"strings"

	"github.com/LNKLEO/OMP/properties"
	"github.com/LNKLEO/OMP/runtime"
	"github.com/LNKLEO/OMP/template"
)

const (
	// Fallback to native command
	NativeFallback properties.Property = "native_fallback"
	// Override the built-in status formats
	StatusFormats properties.Property = "status_formats"
)

// ScmStatus represents part of the status of a repository
type ScmStatus struct {
	Formats    map[string]string
	Unmerged   int
	Deleted    int
	Added      int
	Modified   int
	Moved      int
	Conflicted int
	Untracked  int
	Clean      int
	Missing    int
	Ignored    int
}

func (s *ScmStatus) Changed() bool {
	return s.Unmerged > 0 ||
		s.Added > 0 ||
		s.Deleted > 0 ||
		s.Modified > 0 ||
		s.Moved > 0 ||
		s.Conflicted > 0 ||
		s.Untracked > 0 ||
		s.Clean > 0 ||
		s.Missing > 0 ||
		s.Ignored > 0
}

func (s *ScmStatus) String() string {
	var status strings.Builder

	if s.Formats == nil {
		s.Formats = make(map[string]string)
	}

	stringIfValue := func(value int, name, prefix string) {
		if value <= 0 {
			return
		}

		// allow user override for prefix
		if _, ok := s.Formats[name]; ok {
			status.WriteString(fmt.Sprintf(s.Formats[name], value))
			return
		}

		status.WriteString(fmt.Sprintf(" %s%d", prefix, value))
	}

	stringIfValue(s.Untracked, "Untracked", "?")
	stringIfValue(s.Added, "Added", "+")
	stringIfValue(s.Modified, "Modified", "~")
	stringIfValue(s.Deleted, "Deleted", "-")
	stringIfValue(s.Moved, "Moved", ">")
	stringIfValue(s.Unmerged, "Unmerged", "x")
	stringIfValue(s.Conflicted, "Conflicted", "!")
	stringIfValue(s.Missing, "Missing", "!")
	stringIfValue(s.Clean, "Clean", "=")
	stringIfValue(s.Ignored, "Ignored", "Ø")

	return strings.TrimSpace(status.String())
}

type scm struct {
	base

	Dir             string
	RepoName        string
	mainSCMDir      string
	scmDir          string
	repoRootDir     string
	command         string
	IsWslSharedPath bool
	CommandMissing  bool
	nativeFallback  bool
}

const (
	// BranchTemplate allows to specify a template for the branch name
	BranchTemplate properties.Property = "branch_template"
)

func (s *scm) formatBranch(branch string) string {
	mappedBranches := s.props.GetKeyValueMap(MappedBranches, make(map[string]string))
	for key, value := range mappedBranches {
		matchSubFolders := strings.HasSuffix(key, "*")

		if matchSubFolders && len(key) > 1 {
			key = key[0 : len(key)-1] // remove trailing /* or \*
		}

		if !strings.HasPrefix(branch, key) {
			continue
		}

		branch = strings.Replace(branch, key, value, 1)
		break
	}

	branchTemplate := s.props.GetString(BranchTemplate, "")
	if len(branchTemplate) == 0 {
		return branch
	}

	tmpl := &template.Text{
		Template: branchTemplate,
		Context:  struct{ Branch string }{Branch: branch},
	}

	text, err := tmpl.Render()
	if err != nil {
		return branch
	}

	return text
}

func (s *scm) FileContents(folder, file string) string {
	return strings.Trim(s.env.FileContent(folder+"/"+file), " \r\n")
}

func (s *scm) convertToWindowsPath(path string) string {
	// only convert when in Windows, or when in a WSL shared folder and not using the native fallback
	if s.env.GOOS() == runtime.WINDOWS || (s.IsWslSharedPath && !s.nativeFallback) {
		return s.env.ConvertToWindowsPath(path)
	}

	return path
}

func (s *scm) convertToLinuxPath(path string) string {
	if !s.IsWslSharedPath {
		return path
	}

	return s.env.ConvertToLinuxPath(path)
}

func (s *scm) hasCommand(command string) bool {
	if len(s.command) > 0 {
		return true
	}

	// when in a WSL shared folder, we must use command.exe and convert paths accordingly
	// for worktrees, stashes, and path to work, except when native_fallback is set
	s.IsWslSharedPath = s.env.InWSLSharedDrive()
	if s.env.GOOS() == runtime.WINDOWS || s.IsWslSharedPath {
		command += ".exe"
	}

	if s.env.HasCommand(command) {
		s.command = command
		return true
	}

	s.CommandMissing = true

	// only use the native fallback when set by the user
	if s.IsWslSharedPath && s.props.GetBool(NativeFallback, false) {
		command = strings.TrimSuffix(command, ".exe")
		if s.env.HasCommand(command) {
			s.command = command
			s.nativeFallback = true
			return true
		}
	}

	return false
}
