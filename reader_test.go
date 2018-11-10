package pbzip2

import (
	"bytes"
	"io/ioutil"
	"testing"
)

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
