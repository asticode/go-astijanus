package astijanus

import (
	"flag"

	astihttp "github.com/asticode/go-astitools/http"
)

// Flags
var (
	Addr = flag.String("janus-addr", "", "the Janus addr")
)

// Configuration represents the lib's configuration
type Configuration struct {
	Addr   string `toml:"addr"`
	Sender astihttp.SenderOptions
}

// FlagConfig generates a Configuration based on flags
func FlagConfig() Configuration {
	return Configuration{
		Addr: *Addr,
	}
}
