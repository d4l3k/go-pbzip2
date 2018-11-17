package pbzip2

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// WriterConfig stores configuration options for a pbzip2 writer
type WriterConfig struct {
	Level int
}

type pbzip2Writer struct {
	io.WriteCloser
	cmd *exec.Cmd
}

// NewWriter creates a new pbzip2 writer. This will return an error if pbzip2 is
// not present on the system.
func NewWriter(w io.Writer) (io.WriteCloser, error) {
	defaultConf := newDefaultWriterConfig()

	return newPbzip2Writer(w, defaultConf)
}

// NewWriterConfig creates a new pbzip2 writer with configuration options.
// This will return an error if pbzip2 is not present on the system.
func NewWriterConfig(w io.Writer, conf *WriterConfig) (io.WriteCloser, error) {
	if conf == nil {
		return NewWriter(w)
	}

	return newPbzip2Writer(w, conf)
}

func newPbzip2Writer(w io.Writer, conf *WriterConfig) (io.WriteCloser, error) {
	if !hasPBZip2() {
		return nil, errors.New("missing pbzip2 from system")
	}

	err := conf.validate()
	if err != nil {
		return nil, err
	}

	lvl := fmt.Sprintf("-%d", conf.Level)
	cmd := exec.Command(Command, "-z", lvl)
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

func newDefaultWriterConfig() *WriterConfig {
	return &WriterConfig{
		Level: DefaultCompression,
	}
}

func (wc *WriterConfig) validate() error {
	if wc.Level == 0 {
		wc.Level = DefaultCompression
	}
	if wc.Level < BestSpeed || wc.Level > BestCompression {
		return fmt.Errorf("invalid compression level: %d", wc.Level)
	}

	return nil
}

func (p pbzip2Writer) Close() error {
	if err := p.WriteCloser.Close(); err != nil {
		return err
	}
	return p.cmd.Wait()
}
