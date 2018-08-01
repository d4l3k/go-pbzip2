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
