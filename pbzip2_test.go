package pbzip2

import (
	"bytes"
	"compress/bzip2"
	"crypto/rand"
	"io"
	"io/ioutil"
	"strconv"
	"testing"
)

var targetLen int64 = 10000000

func TestHasPBZip2(t *testing.T) {
	if !hasPBZip2() {
		t.Errorf("should have pbzip2")
	}
}

func testString(t testing.TB) string {
	return testStringN(t, targetLen)
}

func testStringN(t testing.TB, targetLen int64) string {
	var buf bytes.Buffer
	if _, err := io.CopyN(&buf, rand.Reader, targetLen); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if len(out) != int(targetLen) {
		t.Fatal("incorrect testString length")
	}
	return out
}

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	writer, err := NewWriter(&buf)
	if err != nil {
		t.Fatal(err)
	}
	in := testString(t)
	if _, err := writer.Write([]byte(in)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	reader := bzip2.NewReader(&buf)
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	out := string(body)
	if out != in {
		t.Errorf("%q != %q", out, in)
	}
}

func TestReader(t *testing.T) {
	var buf bytes.Buffer
	writer, err := NewWriter(&buf)
	if err != nil {
		t.Fatal(err)
	}
	in := testString(t)
	if _, err := writer.Write([]byte(in)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	reader, err := NewReader(&buf)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	out := string(body)
	if out != in {
		t.Errorf("%q != %q", out, in)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
}

var benchSizes = []int64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000}

func BenchmarkPBZip2Read(b *testing.B) {
	for _, size := range benchSizes {
		size := size
		b.Run(strconv.Itoa(int(size))+"B", func(b *testing.B) {
			var buf bytes.Buffer
			writer, err := NewWriter(&buf)
			if err != nil {
				b.Fatal(err)
			}
			in := testStringN(b, size)
			if _, err := writer.Write([]byte(in)); err != nil {
				b.Fatal(err)
			}
			if err := writer.Close(); err != nil {
				b.Fatal(err)
			}
			compressed := buf.Bytes()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				reader, err := NewReader(bytes.NewReader(compressed))
				if err != nil {
					b.Fatal(err)
				}
				if _, err := ioutil.ReadAll(reader); err != nil {
					b.Fatal(err)
				}
				if err := reader.Close(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkBZip2Read(b *testing.B) {
	for _, size := range benchSizes {
		size := size
		b.Run(strconv.Itoa(int(size))+"B", func(b *testing.B) {
			var buf bytes.Buffer
			writer, err := NewWriter(&buf)
			if err != nil {
				b.Fatal(err)
			}
			in := testStringN(b, size)
			if _, err := writer.Write([]byte(in)); err != nil {
				b.Fatal(err)
			}
			if err := writer.Close(); err != nil {
				b.Fatal(err)
			}
			compressed := buf.Bytes()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				reader := bzip2.NewReader(bytes.NewReader(compressed))
				if _, err := ioutil.ReadAll(reader); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
