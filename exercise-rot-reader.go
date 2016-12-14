package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
  r io.Reader
}

func rot13(b byte) byte {
	result := b + 13
	if (b > 77 && b < 91) || (b > 109 && b < 123) {
		result -= 26
	}
	return result
}

func (reader rot13Reader) Read(bytes []byte) (int, error) {
	n, err := reader.r.Read(bytes)
	for i := range bytes {
		bytes[i] = rot13(bytes[i])
	}
	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
