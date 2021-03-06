package prefixedio

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	var fooBuf, barBuf bytes.Buffer

	original := bytes.NewReader([]byte(strings.TrimSpace(testInput)))
	r, err := NewReader(original)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	pFoo, err := r.Prefix("foo: ")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	_, err = r.Prefix("foo: ")
	if err == nil {
		t.Fatalf("expected prefix already registered error")
	}
	pBar, err := r.Prefix("bar: ")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	doneCh := make(chan struct{}, 2)
	go func() {
		if _, err := io.Copy(&fooBuf, pFoo); err != nil {
			t.Fatalf("err: %s", err)
		}
		doneCh <- struct{}{}
	}()
	go func() {
		if _, err := io.Copy(&barBuf, pBar); err != nil {
			t.Fatalf("err: %s", err)
		}
		doneCh <- struct{}{}
	}()

	// Wait for all the reads to be done
	<-doneCh
	<-doneCh

	if fooBuf.String() != testInputFoo {
		t.Fatalf("bad: %s", fooBuf.String())
	}
	if barBuf.String() != testInputBar {
		t.Fatalf("bad: %s", barBuf.String())
	}
}

const testInput = `
foo: 1
bar: 2
bar: 3
42
foo: 4
foo: 5
bar: 6
`

const testInputFoo = `1
4
5
`

const testInputBar = `2
3
42
6`
