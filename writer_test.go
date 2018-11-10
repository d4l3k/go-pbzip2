package pbzip2

import (
	"bytes"
	"compress/bzip2"
	"io/ioutil"
	"testing"
)

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

func TestWriterConfig(t *testing.T) {
	var buf bytes.Buffer
	conf := &WriterConfig{Level: 6}
	writer, err := NewWriterConfig(&buf, conf)
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

func TestWriterConfigNil(t *testing.T) {
	var buf bytes.Buffer
	writer, err := NewWriterConfig(&buf, nil)
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

func TestNewDefaultWriterConfig(t *testing.T) {
	defaultConf := newDefaultWriterConfig()

	if defaultConf.Level != DefaultCompression {
		t.Errorf("expected default config level %d, got: %d",
			DefaultCompression, defaultConf.Level)
	}
}

func TestWriterConfigValidate(t *testing.T) {
	// Level of 0 gets set to default
	conf1 := &WriterConfig{Level: 0}
	err := conf1.validate()
	if err != nil {
		t.Fatal(err)
	}
	if conf1.Level != DefaultCompression {
		t.Errorf("expected default config level %d, got: %d",
			DefaultCompression, conf1.Level)
	}

	// Level out of lower range
	conf2 := &WriterConfig{Level: -1}
	err = conf2.validate()
	if err == nil {
		t.Error("exptected error with writer config Level of -1")
	}

	// Level out of upper range
	conf3 := &WriterConfig{Level: 10}
	err = conf3.validate()
	if err == nil {
		t.Error("exptected error with writer config Level of 10")
	}
}
