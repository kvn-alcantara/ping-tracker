package factory

import (
	"runtime"

	"github.com/kvn-alcantara/ping-tracker/internal/repository"
)

// NewPinger creates and returns a platform-specific implementation of the Pinger interface.
func NewPinger() repository.Pinger {
	if runtime.GOOS == "linux" {
		return repository.NewHTTPPinger()
	}
	return repository.NewProBingPinger()
}
