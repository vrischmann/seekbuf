package seekbuf

import (
	"fmt"
	"io"
	"os"
)

type Buffer struct {
	data []byte
	pos  int
}

// Bytes returns a slice holding the data from the current position up to the end.
func (b *Buffer) Bytes() []byte {
	return b.data[b.pos:]
}

// Read reads the next len(p) bytes from the buffer starting from the current position.
func (b *Buffer) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	if len(b.data[b.pos:]) <= 0 {
		return 0, io.EOF
	}

	n := copy(p, b.data[b.pos:])
	b.pos += n

	return n, nil
}

// Write the contents of p in the buffer at the current position, growing the buffer if needed.
func (b *Buffer) Write(p []byte) (int, error) {
	n := copy(b.data[b.pos:], p)
	b.data = append(b.data, p[n:]...)
	b.pos += len(p)

	return len(p), nil
}

// Seek sets the offset for the next Read or Write on the buffer to offset, interpreted according to whence:
// 0 means relative to the origin of the buffer, 1 means relative to the current offset, and 2 means relative to the end.
// It returns the new offset and an error, if any.
func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	o := int(offset)
	switch whence {
	case os.SEEK_CUR:
		if o > 0 && b.pos+o >= len(b.data) {
			return -1, fmt.Errorf("invalid offset %d", offset)
		}
		b.pos += o

	case os.SEEK_SET:
		if o > 0 && o >= len(b.data) {
			return -1, fmt.Errorf("invalid offset %d", offset)
		}
		b.pos = o

	case os.SEEK_END:
		if len(b.data)+o < 0 {
			return -1, fmt.Errorf("invalid offset %d", offset)
		}
		b.pos = len(b.data) + o

	default:
		return -1, fmt.Errorf("invalid whence %d", whence)
	}

	return int64(b.pos), nil
}

var _ io.ReadWriteSeeker = (*Buffer)(nil)
