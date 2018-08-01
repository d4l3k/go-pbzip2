package pbzip2

import "io"

type bzip2Closer struct {
	io.Reader
}

func (bzip2Closer) Close() error {
	return nil
}

var _ io.ReadCloser = bzip2Closer{}
