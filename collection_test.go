package bcs_test

import (
	"errors"
	"io"
	"math"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/iotaledger/bcs-go"
)

func TestArrayCodec(t *testing.T) {
	bcs.TestCodecAndBytes(t, []int64{42, 43}, []byte{0x2, 0x2A, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2B, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	bcs.TestCodecAndBytes(t, []int8{42, 43}, []byte{0x2, 0x2A, 0x2B})
	bcs.TestCodecAndBytes(t, []int8(nil), []byte{0x0})
	bcs.TestCodecAndBytes(t, []uint8{42, 43}, []byte{0x2, 0x2A, 0x2B})
	bcs.TestCodecAndBytes(t, []int64(nil), []byte{0x0})

	bcs.TestCodecAndBytes(t, []*int16{lo.ToPtr[int16](1), lo.ToPtr[int16](2), lo.ToPtr[int16](3)}, []byte{0x3, 0x1, 0x0, 0x2, 0x0, 0x3, 0x0})
	bcs.TestCodecAndBytes(t, []*byte{lo.ToPtr[byte](42), lo.ToPtr[byte](43)}, []byte{0x2, 0x2A, 0x2B})
	bcs.TestCodecAndBytes(t, []*int8{lo.ToPtr[int8](42), lo.ToPtr[int8](43)}, []byte{0x2, 0x2A, 0x2B})

	bcs.TestCodecAndBytes(t, []BasicWithCustomCodec{"a", "b"}, []byte{0x2, 0x1, 0x2, 0x3, 0x1, 0x61, 0x1, 0x2, 0x3, 0x1, 0x62})
	bcs.TestCodecAndBytes(t, []*BasicWithCustomCodec{lo.ToPtr[BasicWithCustomCodec]("a"), lo.ToPtr[BasicWithCustomCodec]("b")}, []byte{0x2, 0x1, 0x2, 0x3, 0x1, 0x61, 0x1, 0x2, 0x3, 0x1, 0x62})
	bcs.TestCodecAndBytes(t, []BasicWithCustomPtrCodec{"a", "b"}, []byte{0x2, 0x1, 0x2, 0x3, 0x1, 0x61, 0x1, 0x2, 0x3, 0x1, 0x62})

	bcs.TestCodecAndBytes(t, [3]int64{42, 43, 44}, []byte{0x2a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	bcs.TestCodecAndBytes(t, [3]byte{42, 43, 44}, []byte{0x2a, 0x2b, 0x2c})

	bcs.TestCodecAndBytes(t, []string{"aaa", "bbb"}, []byte{0x2, 0x3, 0x61, 0x61, 0x61, 0x3, 0x62, 0x62, 0x62})
	bcs.TestCodecAndBytes(t, [][]int16{{1, 2}, {3, 4, 5}}, []byte{0x2, 0x2, 0x1, 0x0, 0x2, 0x0, 0x3, 0x3, 0x0, 0x4, 0x0, 0x5, 0x0})
}

func TestDecodeMalformedRegularSlice(t *testing.T) {
	e := bcs.NewBytesEncoder()

	const elemsCount100Billions = 100_000_000_000
	e.WriteLen(elemsCount100Billions)
	e.Encode(BasicWithCustomCodec("hello"))
	e.Encode(BasicWithCustomCodec("world"))

	_, err := bcs.Unmarshal[[]BasicWithCustomCodec](e.Bytes())
	require.Error(t, err)
	require.True(t, errors.Is(err, io.EOF))
}

func TestDecodeMalformedBytesSlice(t *testing.T) {
	e := bcs.NewBytesEncoder()

	const ramSize1000GB = 1000 * 1024 * 1024 * 1024
	e.WriteLen(ramSize1000GB)
	e.WriteByte(1)
	e.WriteByte(2)
	e.WriteByte(3)

	_, err := bcs.Unmarshal[[]byte](e.Bytes())
	require.Error(t, err)
	require.True(t, errors.Is(err, io.EOF))
}

func TestMapCodec(t *testing.T) {
	intMapEnc := []byte{0x3, 0x0, 0x0, 0x0, 0x3, 0x0, 0x1, 0xfd, 0xff, 0x1}
	bcs.TestCodecAndBytes(t, map[int16]bool{-3: true, 0: false, 3: true}, intMapEnc)
	bcs.TestCodecAndBytes(t, map[int16]bool{3: true, 0: false, -3: true}, intMapEnc)
	bcs.TestCodecAndBytes(t, map[int16]bool{}, []byte{0x0})

	uintMapEnc := []byte{0x3, 0x1, 0x0, 0x0, 0x2, 0x0, 0x1, 0x3, 0x0, 0x1}
	bcs.TestCodecAndBytes(t, map[uint16]bool{3: true, 1: false, 2: true}, uintMapEnc)
	bcs.TestCodecAndBytes(t, map[uint16]bool{2: true, 1: false, 3: true}, uintMapEnc)
	bcs.TestCodecAndBytes(t, map[uint16]bool{}, []byte{0x0})

	strMapEnc := []byte{0x3, 0x2, 0x61, 0x61, 0x0, 0x2, 0x62, 0x62, 0x1, 0x2, 0x63, 0x63, 0x1}
	bcs.TestCodecAndBytes(t, map[string]bool{"cc": true, "aa": false, "bb": true}, strMapEnc)
	bcs.TestCodecAndBytes(t, map[string]bool{"bb": true, "aa": false, "cc": true}, strMapEnc)

	intMapOfMapsEnc := []byte{0x2, 0x1, 0x0, 0x2, 0x2, 0x0, 0x1, 0x3, 0x0, 0x0, 0x2, 0x0, 0x1, 0x1, 0x0, 0x1}
	bcs.TestCodecAndBytes(t, map[int16]map[int16]bool{1: {2: true, 3: false}, 2: {1: true}}, intMapOfMapsEnc)

	customMapEnc := []byte{0x2, 0x1, 0x2, 0x3, 0x2, 0x61, 0x61, 0x1, 0x2, 0x3, 0x1, 0x2, 0x3, 0x2, 0x62, 0x62, 0x1, 0x2, 0x3}
	bcs.TestCodecAndBytes(t, map[BasicWithCustomCodec]WithCustomCodec{"bb": {}, "aa": {}}, customMapEnc)
	bcs.TestCodecAndBytes(t, map[BasicWithCustomCodec]WithCustomCodec{"aa": {}, "bb": {}}, customMapEnc)
	bcs.TestCodecAndBytes(t, map[BasicWithCustomCodec]*WithCustomCodec{"bb": {}, "aa": {}}, customMapEnc)
	bcs.TestCodecAndBytes(t, map[BasicWithCustomCodec]*WithCustomCodec{"aa": {}, "bb": {}}, customMapEnc)
	bcs.TestCodecAndBytes(t, map[BasicWithCustomPtrCodec]*WithCustomCodec{"aa": {}, "bb": {}}, customMapEnc)

	customMapEnc = []byte{0x2, 0x1, 0x2, 0x3, 0x2, 0x61, 0x61, 0x1, 0x2, 0x3, 0x2, 0x63, 0x63, 0x1, 0x2, 0x3, 0x2, 0x62, 0x62, 0x1, 0x2, 0x3, 0x2, 0x64, 0x64}
	bcs.TestCodecAndBytes(t, map[BasicWithCustomPtrCodec]BasicWithCustomPtrCodec{"aa": "cc", "bb": "dd"}, customMapEnc)
}

func TestCollectionSizeCodec(t *testing.T) {
	testSizeCodec(t, 0, []byte{0x0})
	testSizeCodec(t, 1, []byte{0x1})
	testSizeCodec(t, 127, []byte{0x7F})
	testSizeCodec(t, 128, []byte{0x80, 0x1})
	testSizeCodec(t, 16383, []byte{0xFF, 0x7F})
	testSizeCodec(t, 16384, []byte{0x80, 0x80, 0x1})
	testSizeCodec(t, 2097151, []byte{0xFF, 0xFF, 0x7F})
	testSizeCodec(t, 2097152, []byte{0x80, 0x80, 0x80, 0x1})
	testSizeCodec(t, 268435455, []byte{0xFF, 0xFF, 0xFF, 0x7F})
	testSizeCodec(t, 268435456, []byte{0x80, 0x80, 0x80, 0x80, 0x1})
	testSizeCodec(t, 2147483647, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x7})
	testSizeCodec(t, 4294967295, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xF})
	testSizeCodec(t, math.MaxUint32, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xF})
}

func testSizeCodec(t *testing.T, v int, expectedEnc []byte) {
	e := bcs.NewBytesEncoder()
	e.WriteLen(v)
	require.Equal(t, expectedEnc, e.Bytes())
}
