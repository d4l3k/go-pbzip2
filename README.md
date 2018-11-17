# go-pbzip2 [![GoDoc](https://godoc.org/github.com/d4l3k/go-pbzip2?status.svg)](https://godoc.org/github.com/d4l3k/go-pbzip2) [![Build Status](https://travis-ci.com/d4l3k/go-pbzip2.svg?branch=master)](https://travis-ci.com/d4l3k/go-pbzip2)

A go library to use pbzip2 for faster bzip2 operations than the stdlib. Supports
both compression and decompression.

## Benchmarks

For small amounts of bzip2 data, using the standard bzip2 library is probably
faster. Once you start decompressing data beyond 1MB, pbzip2 will be a lot
faster with roughly linear performance to the number of cores. These benchmarks
are run with `%d` random bytes compressed with pbzip2. Machine was an i7-8550
processor with 4 cores and 8 threads.

```
goos: linux
goarch: amd64
pkg: github.com/d4l3k/go-pbzip2
BenchmarkPBZip2Read/1B-8                    5000           3509845 ns/op
BenchmarkPBZip2Read/10B-8                   5000           3306304 ns/op
BenchmarkPBZip2Read/100B-8                  5000           3315698 ns/op
BenchmarkPBZip2Read/1000B-8                 5000           3663066 ns/op
BenchmarkPBZip2Read/10000B-8                3000           4580623 ns/op
BenchmarkPBZip2Read/100000B-8               2000          11377153 ns/op
BenchmarkPBZip2Read/1000000B-8               200          66438579 ns/op
BenchmarkPBZip2Read/10000000B-8               50         329762754 ns/op
BenchmarkPBZip2Read/100000000B-8               5        2416747246 ns/op
BenchmarkBZip2Read/1B-8                    30000            388585 ns/op
BenchmarkBZip2Read/10B-8                   30000            401892 ns/op
BenchmarkBZip2Read/100B-8                  30000            492336 ns/op
BenchmarkBZip2Read/1000B-8                 20000            743518 ns/op
BenchmarkBZip2Read/10000B-8                10000           1981048 ns/op
BenchmarkBZip2Read/100000B-8                2000           9123679 ns/op
BenchmarkBZip2Read/1000000B-8                200          76655666 ns/op
BenchmarkBZip2Read/10000000B-8                20         745520450 ns/op
BenchmarkBZip2Read/100000000B-8                2        8047164946 ns/op
```

## Usage

### Reader

```go
// If pbzip2 is not present on the system, stdlib bzip2.Reader is used instead.
func NewReader(r io.Reader) (io.ReadCloser, error) { ... }
```

### Writer

```go
// NewWriter will use the default Level of 9
func NewWriter(w io.Writer) (io.WriteCloser, error) { ... }

type WriterConfig struct {
    Level int
}

func NewWriterConfig(w io.Writer, conf *WriterConfig) (io.WriteCloser, error) { ... }
```
