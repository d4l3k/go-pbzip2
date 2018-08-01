// pbzip2 is a wrapper around the pbzip2 command. If it's present on the system,
// this library will use it for higher performance when decompressing bzip2
// data.
package pbzip2

import (
	"compress/bzip2"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

const Command = "pbzip2"

func hasPBZip2() bool {
	cmd := exec.Command(Command, "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// NewReader creates a new pbzip2 reader. This will print a warning if pbzip2 is
// not present on the system and return a stdlib bzip2.Reader instead.
func NewReader(r io.Reader) (io.ReadCloser, error) {
	if !hasPBZip2() {
		warn()
		return bzip2Closer{bzip2.NewReader(r)}, nil
	}

	cmd := exec.Command(Command, "-d")
	cmd.Stderr = os.Stderr
	cmd.Stdin = r

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return pbzip2Reader{
		ReadCloser: out,
		cmd:        cmd,
	}, nil
}

type pbzip2Reader struct {
	io.ReadCloser
	cmd *exec.Cmd
}

func (p pbzip2Reader) Close() error {
	if err := p.ReadCloser.Close(); err != nil {
		return err
	}
	return p.cmd.Wait()
}

// NewWriter creates a new pbzip2 writer. This will return an error if pbzip2 is
// not present on the system.
func NewWriter(w io.Writer) (io.WriteCloser, error) {
	if !hasPBZip2() {
		return nil, errors.Errorf("missing pbzip2 from system")
	}

	cmd := exec.Command(Command, "-z")
	cmd.Stderr = os.Stderr
	cmd.Stdout = w

	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return pbzip2Writer{
		WriteCloser: in,
		cmd:         cmd,
	}, nil
}

type pbzip2Writer struct {
	io.WriteCloser
	cmd *exec.Cmd
}

func (p pbzip2Writer) Close() error {
	if err := p.WriteCloser.Close(); err != nil {
		return err
	}
	return p.cmd.Wait()
}
