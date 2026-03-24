// core/config.go v1
package core

import "time"

type Config struct {
	Timeout             time.Duration
	BitLenLimit         int
	BitLenWarnThreshold int
}

func DefaultConfig() Config {
	return Config{}
}

// core/config.go v1
