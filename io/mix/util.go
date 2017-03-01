package mix

import (
	"errors"
	"io"
)

// ReadAtSeeker embeds io.ReadSeeker and io.ReaderAt
// To be a ReadAtSeeker, a type must implement Read, Seek, and ReadAt
type ReadAtSeeker interface {
	io.ReadSeeker
	io.ReaderAt
}

// Reads in a fixed length UTF8 string.
func ReadFixedUTF8String(rdr io.Reader, size int) (string, error) {
	out := make([]byte, size)

	n, err := io.ReadFull(rdr, out)
	if err != nil {
		return "", err
	} else if n != size {
		return "", errors.New("ReadFixedUTF8String: io.ReadFull couldn't read the full size given.")
	}

	// []byte => string type conversion seems to assume UTF8.
	// Though I can't find any document describing this.
	return string(out), nil
}
