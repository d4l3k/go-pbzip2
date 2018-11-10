// Package pbzip2 is a wrapper around the pbzip2 command. If it's present on
// the system, this library will use it for higher performance when
// decompressing bzip2 data.
package pbzip2

import (
	"os/exec"
)

const (
	// Command is the pbzip2 command
	Command = "pbzip2"

	// BestSpeed represents the smallest block size (100k)
	BestSpeed = 1
	// BestCompression represents the largest block size (900k)
	BestCompression = 9
	// DefaultCompression represents the default block size (900k)
	DefaultCompression = 9
)

func hasPBZip2() bool {
	cmd := exec.Command(Command, "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
