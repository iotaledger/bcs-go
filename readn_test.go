package bcs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadN(t *testing.T) {
	testReadN(t, 1, 1)
	testReadN(t, 1, 0)
	testReadN(t, 0, 0)
	testReadN(t, 100, 10)
	testReadN(t, 100, 50)
	testReadN(t, 100, 100)
	testReadN(t, maxReadNBufferSize, 0)
	testReadN(t, maxReadNBufferSize, 10)
	testReadN(t, maxReadNBufferSize, maxReadNBufferSize-1)
	testReadN(t, maxReadNBufferSize, maxReadNBufferSize)
	testReadN(t, maxReadNBufferSize+100, 0)
	testReadN(t, maxReadNBufferSize+100, 10)
	testReadN(t, maxReadNBufferSize+100, maxReadNBufferSize-1)
	testReadN(t, maxReadNBufferSize+100, maxReadNBufferSize)
	testReadN(t, maxReadNBufferSize*3, 0)
	testReadN(t, maxReadNBufferSize*3, 10)
	testReadN(t, maxReadNBufferSize*3, maxReadNBufferSize-1)
	testReadN(t, maxReadNBufferSize*3, maxReadNBufferSize)
	testReadN(t, maxReadNBufferSize*3, maxReadNBufferSize+1)
	testReadN(t, maxReadNBufferSize*3, maxReadNBufferSize*2-1)
	testReadN(t, maxReadNBufferSize*3, maxReadNBufferSize*2)
	testReadN(t, maxReadNBufferSize*3, maxReadNBufferSize*3-1)
	testReadN(t, maxReadNBufferSize*3, maxReadNBufferSize*3)
}

func testReadN(t *testing.T, dataSize, bytesToRead int) {
	b := make([]byte, dataSize)
	for i := 0; i < dataSize; i++ {
		b[i] = byte(i) + 1
	}

	r := bytes.NewReader(b)
	d := NewDecoder(r)

	res, err := d.ReadN(bytesToRead)
	require.NoError(t, err)
	require.Equal(t, b[:bytesToRead], res)
}
