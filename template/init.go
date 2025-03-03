package template

import (
	"sync"

	"github.com/LNKLEO/OMP/maps"
	"github.com/LNKLEO/OMP/runtime"
)

const (
	// Errors to show when the template handling fails
	InvalidTemplate   = "invalid template text"
	IncorrectTemplate = "unable to create text based on template"

	globalRef = ".$"
)

var (
	shell       string
	env         runtime.Environment
	knownFields *maps.Concurrent
)

func Init(environment runtime.Environment, vars maps.Simple) {
	env = environment
	shell = env.Shell()
	knownFields = maps.NewConcurrent()

	renderPool = sync.Pool{
		New: func() any {
			return newTextPoolObject()
		},
	}

	if Cache != nil {
		return
	}

	loadCache(vars)
}
