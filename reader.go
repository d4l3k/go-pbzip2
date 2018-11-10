package pbzip2

import (
	"compress/bzip2"
	"io"
	"os"
	"os/exec"
)

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
