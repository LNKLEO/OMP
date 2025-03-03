package segments

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/LNKLEO/OMP/log"
)

const (
	PoshGitEnv = "OMP_GIT_STATUS"
)

type PoshGit struct {
	Index        *PoshGitStatus `json:"Index"`
	Working      *PoshGitStatus `json:"Working"`
	RepoName     string         `json:"RepoName"`
	Branch       string         `json:"Branch"`
	GitDir       string         `json:"GitDir"`
	Upstream     string         `json:"Upstream"`
	StashCount   int            `json:"StashCount"`
	AheadBy      int            `json:"AheadBy"`
	BehindBy     int            `json:"BehindBy"`
	HasWorking   bool           `json:"HasWorking"`
	HasIndex     bool           `json:"HasIndex"`
	HasUntracked bool           `json:"HasUntracked"`
}

type PoshGitStatus struct {
	Added    []string `json:"Added"`
	Modified []string `json:"Modified"`
	Deleted  []string `json:"Deleted"`
	Unmerged []string `json:"Unmerged"`
}

func (s *GitStatus) parseGitStatus(p *PoshGitStatus) {
	if p == nil {
		return
	}

	s.Added = len(p.Added)
	s.Deleted = len(p.Deleted)
	s.Modified = len(p.Modified)
	s.Unmerged = len(p.Unmerged)
}

func (g *Git) hasGitStatus() bool {
	envStatus := g.env.Getenv(PoshGitEnv)
	if len(envStatus) == 0 {
		log.Error(fmt.Errorf("%s environment variable not set, do you have the posh-git module installed?", PoshGitEnv))
		return false
	}

	var git PoshGit
	err := json.Unmarshal([]byte(envStatus), &git)
	if err != nil {
		log.Error(err)
		return false
	}

	g.setDir(git.GitDir)
	g.Working = &GitStatus{}
	g.Working.parseGitStatus(git.Working)
	g.Staging = &GitStatus{}
	g.Staging.parseGitStatus(git.Index)
	g.HEAD = g.parseGitHEAD(git.Branch)
	g.stashCount = git.StashCount
	g.Ahead = git.AheadBy
	g.Behind = git.BehindBy
	g.UpstreamGone = len(git.Upstream) == 0
	g.Upstream = git.Upstream

	g.setBranchStatus()

	if len(g.Upstream) != 0 && g.props.GetBool(FetchUpstreamIcon, false) {
		g.UpstreamIcon = g.getUpstreamIcon()
	}

	g.PoshGit = true
	return true
}

func (g *Git) parseGitHEAD(head string) string {
	// commit
	if strings.HasSuffix(head, "...)") {
		head = strings.TrimLeft(head, "(")
		head = strings.TrimRight(head, ".)")
		return fmt.Sprintf("%s%s", g.props.GetString(CommitIcon, "\uF417"), head)
	}
	// tag
	if strings.HasPrefix(head, "(") {
		head = strings.TrimLeft(head, "(")
		head = strings.TrimRight(head, ")")
		return fmt.Sprintf("%s%s", g.props.GetString(TagIcon, "\uF412"), head)
	}
	// regular branch
	return fmt.Sprintf("%s%s", g.props.GetString(BranchIcon, "\uE0A0"), g.formatBranch(head))
}
