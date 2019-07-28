package vlq

import (
	"bytes"
	"testing"
)

func assertRvlqEncodedReversed(t *testing.T, value int, expectedNonReversedEncoded []byte) {
	byteCount := len(expectedNonReversedEncoded)
	inflatedByteCount := byteCount + 10
	expectedEncoded := make([]byte, inflatedByteCount)
	for i, v := range expectedNonReversedEncoded {
		expectedEncoded[inflatedByteCount-i-1] = v
	}
	actualEncoded := make([]byte, inflatedByteCount)
	vlq := Rvlq(value)

	bytesUsed, err := vlq.EncodeReversedTo(actualEncoded)
	if err != nil {
		t.Error(err)
		return
	}
	if bytesUsed != byteCount {
		t.Errorf("Expected to use %v bytes but actually used %v", byteCount, bytesUsed)
		return
	}
	if !bytes.Equal(expectedEncoded, actualEncoded) {
		t.Errorf("Expected %v but got %v", expectedEncoded, actualEncoded)
		return
	}
	decoded, bytesUsed, err := DecodeRvlqReversedFrom(actualEncoded)
	if err != nil {
		t.Error(err)
		return
	}
	if bytesUsed != byteCount {
		t.Errorf("Expected to use %v bytes but actually used %v", byteCount, bytesUsed)
		return
	}
	if decoded != vlq {
		t.Errorf("Expected decoded value %v but got %v", vlq, decoded)
	}
}

func assertRvlqEncoded(t *testing.T, value int, expectedEncoded []byte) {
	vlq := Rvlq(value)
	byteCount := vlq.EncodedSize()
	actualEncoded := make([]byte, byteCount)
	bytesUsed, err := vlq.EncodeTo(actualEncoded)
	if err != nil {
		t.Error(err)
		return
	}
	if bytesUsed != byteCount {
		t.Errorf("Expected to use %v bytes but actually used %v", byteCount, bytesUsed)
		return
	}
	if !bytes.Equal(expectedEncoded, actualEncoded) {
		t.Errorf("Expected %v but got %v", expectedEncoded, actualEncoded)
	}
	decoded, bytesUsed, isComplete := DecodeRvlqFrom(actualEncoded)
	if bytesUsed != byteCount {
		t.Errorf("Expected to use %v bytes but actually used %v", byteCount, bytesUsed)
		return
	}
	if !isComplete {
		t.Errorf("Expected decoding to be complete")
		return
	}
	if decoded != vlq {
		t.Errorf("Expected decoded value %v but got %v", vlq, decoded)
		return
	}

	assertRvlqEncodedReversed(t, value, expectedEncoded)
}

func TestRvlq0(t *testing.T) {
	assertRvlqEncoded(t, 0, []byte{0})
}

func TestRvlq1(t *testing.T) {
	assertRvlqEncoded(t, 1, []byte{1})
}

func TestRvlq7F(t *testing.T) {
	assertRvlqEncoded(t, 0x7f, []byte{0x7f})
}

func TestRvlq80(t *testing.T) {
	assertRvlqEncoded(t, 0x80, []byte{0x81, 0x00})
}

func TestRvlqFF(t *testing.T) {
	assertRvlqEncoded(t, 0xff, []byte{0x81, 0x7f})
}

func TestRvlq2005(t *testing.T) {
	assertRvlqEncoded(t, 0x2005, []byte{0xc0, 0x05})
}

func TestRvlq3FFF(t *testing.T) {
	assertRvlqEncoded(t, 0x3fff, []byte{0xff, 0x7f})
}

func TestRvlq4000(t *testing.T) {
	assertRvlqEncoded(t, 0x4000, []byte{0x81, 0x80, 0x00})
}

func TestRvlqFFFF(t *testing.T) {
	assertRvlqEncoded(t, 0xffff, []byte{0x83, 0xff, 0x7f})
}

func TestRvlq10000(t *testing.T) {
	assertRvlqEncoded(t, 0x10000, []byte{0x84, 0x80, 0x00})
}

func TestRvlq1fffff(t *testing.T) {
	assertRvlqEncoded(t, 0x1fffff, []byte{0xff, 0xff, 0x7f})

}

func TestRvlq200000(t *testing.T) {
	assertRvlqEncoded(t, 0x200000, []byte{0x81, 0x80, 0x80, 0x00})

}
