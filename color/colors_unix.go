//go:build !windows && !darwin

package color

import "github.com/LNKLEO/OMP/runtime"

func GetAccentColor(_ runtime.Environment) (*RGB, error) {
	return nil, &runtime.NotImplemented{}
}
