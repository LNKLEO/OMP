package path

import (
	"os"

	"github.com/LNKLEO/OMP/log"
)

func Home() string {
	home := os.Getenv("HOME")
	defer func() {
		log.Debug(home)
	}()

	if len(home) > 0 {
		return home
	}

	// fallback to older implemenations on Windows
	home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")

	if len(home) == 0 {
		home = os.Getenv("USERPROFILE")
	}

	return home
}
