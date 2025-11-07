package utils

import (
	"bytes"
	"io"
)

func CopyReader(reader io.Reader) (*bytes.Buffer, *bytes.Buffer, error) {

	var buf bytes.Buffer
	var buf2 bytes.Buffer
	tee := io.TeeReader(reader, &buf)

	_, err := io.Copy(&buf2, tee)

	if err != nil {
		return nil, nil, err
	}
	return &buf, &buf2, err
}
