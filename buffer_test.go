package seekbuf_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vrischmann/seekbuf"
)

type readWriteTestCase struct {
	input   string
	output  string
	expWErr string
	expRErr string
}

var readWriteTestCases = []readWriteTestCase{
	{"foobar", "foobar", "", ""},
	{"", "foo", "", "EOF"},
}

func TestReadWrite(t *testing.T) {
	var b seekbuf.Buffer

	n, err := b.Write([]byte("foobar"))
	require.Nil(t, err)
	require.Equal(t, 6, n)

	b.Seek(0, io.SeekStart)

	var data [6]byte
	n, err = b.Read(data[:])
	require.Nil(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, "foobar", string(data[:]))
}

func TestReadWriteBiggerReadBuf(t *testing.T) {
	var b seekbuf.Buffer

	n, err := b.Write([]byte("foo"))
	require.Nil(t, err)
	require.Equal(t, 3, n)

	b.Seek(0, io.SeekStart)

	var data [6]byte
	n, err = b.Read(data[:])
	require.Nil(t, err)
	require.Equal(t, 3, n)
	require.Equal(t, "foo", string(data[:n]))
}

func TestReadWriteEOF(t *testing.T) {
	var b seekbuf.Buffer

	var data [6]byte
	n, err := b.Read(data[:])
	require.NotNil(t, err)
	require.Equal(t, 0, n)
	require.Equal(t, io.EOF, err)
}

func TestReadWriteSeekSet(t *testing.T) {
	var b seekbuf.Buffer

	n, err := b.Write([]byte("foobar"))
	require.Nil(t, err)
	require.Equal(t, 6, n)

	b.Seek(0, io.SeekStart)

	var data [6]byte
	n, err = b.Read(data[:])
	require.Nil(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, "foobar", string(data[:]))

	n2, err := b.Seek(0, io.SeekStart)
	require.Nil(t, err)
	require.Equal(t, int64(0), n2)

	n, err = b.Read(data[:])
	require.Nil(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, "foobar", string(data[:]))
}

func TestReadWriteSeekCur(t *testing.T) {
	var b seekbuf.Buffer

	n, err := b.Write([]byte("foobar"))
	require.Nil(t, err)
	require.Equal(t, 6, n)

	n2, err := b.Seek(3, io.SeekStart)
	require.Nil(t, err)
	require.Equal(t, int64(3), n2)

	require.Equal(t, 3, len(b.Bytes()))
	require.Equal(t, "bar", string(b.Bytes()))

	n2, err = b.Seek(1, io.SeekCurrent)
	require.Nil(t, err)
	require.Equal(t, int64(4), n2)

	require.Equal(t, 2, len(b.Bytes()))
	require.Equal(t, "ar", string(b.Bytes()))
}

func TestReadWriteSeekEnd(t *testing.T) {
	var b seekbuf.Buffer

	n, err := b.Write([]byte("foobar"))
	require.Nil(t, err)
	require.Equal(t, 6, n)

	n2, err := b.Seek(-1, io.SeekEnd)
	require.Nil(t, err)
	require.Equal(t, int64(5), n2)

	require.Equal(t, 1, len(b.Bytes()))
	require.Equal(t, "r", string(b.Bytes()))
}

func TestSeekEmpty(t *testing.T) {
	var b seekbuf.Buffer

	n, err := b.Seek(0, io.SeekStart)
	require.Nil(t, err)
	require.Equal(t, int64(0), n)

	n, err = b.Seek(0, io.SeekCurrent)
	require.Nil(t, err)
	require.Equal(t, int64(0), n)

	n, err = b.Seek(0, io.SeekEnd)
	require.Nil(t, err)
	require.Equal(t, int64(0), n)
}

func TestSeekZeroNotEmpty(t *testing.T) {
	var b seekbuf.Buffer

	b.Write([]byte("foobar"))

	n, err := b.Seek(0, io.SeekStart)
	require.Nil(t, err)
	require.Equal(t, int64(0), n)

	n, err = b.Seek(0, io.SeekCurrent)
	require.Nil(t, err)
	require.Equal(t, int64(0), n)

	n, err = b.Seek(0, io.SeekEnd)
	require.Nil(t, err)
	require.Equal(t, int64(6), n)
}

func TestSeekOverwrite(t *testing.T) {
	var b seekbuf.Buffer

	b.Write([]byte("foobar"))
	b.Seek(0, io.SeekStart)
	b.Write([]byte("baz"))

	b.Seek(0, io.SeekStart)

	require.Equal(t, "bazbar", string(b.Bytes()))
}

func TestWrite(t *testing.T) {
	var b seekbuf.Buffer

	n, err := b.Write([]byte("foobar"))
	require.Nil(t, err)
	require.Equal(t, 6, n)
	n, err = b.Write([]byte("quxbaz"))
	require.Nil(t, err)
	require.Equal(t, 6, n)

	b.Seek(0, io.SeekStart)

	require.Equal(t, "foobarquxbaz", string(b.Bytes()))
}

func TestNew(t *testing.T) {
	data := []byte("foobarquxbaz")
	b := seekbuf.New(data)

	require.Equal(t, "foobarquxbaz", string(b.Bytes()))

	n, err := b.Seek(3, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(3), n)

	var buf [10]byte
	n2, err := b.Read(buf[:])
	require.NoError(t, err)
	require.Equal(t, 9, n2)
	require.Equal(t, "barquxbaz", string(buf[:n2]))
}
