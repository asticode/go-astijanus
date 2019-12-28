package astijanus

import (
	"flag"

	"github.com/asticode/go-astikit"
)

// Flags
var (
	Addr = flag.String("janus-addr", "", "the Janus addr")
)

// Configuration represents the lib's configuration
type Configuration struct {
	Addr   string `toml:"addr"`
	Sender astikit.HTTPSenderOptions
}

// FlagConfig generates a Configuration based on flags
func FlagConfig() Configuration {
	return Configuration{
		Addr: *Addr,
	}
}
