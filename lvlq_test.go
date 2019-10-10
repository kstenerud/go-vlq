package vlq

import (
	"bytes"
	"testing"
)

func assertLvlqEncodedReversed(t *testing.T, value int, expectedNonReversedEncoded []byte) {
	byteCount := len(expectedNonReversedEncoded)
	inflatedByteCount := byteCount + 10
	expectedEncoded := make([]byte, inflatedByteCount)
	for i, v := range expectedNonReversedEncoded {
		expectedEncoded[inflatedByteCount-i-1] = v
	}
	actualEncoded := make([]byte, inflatedByteCount)
	vlq := Lvlq(value)

	bytesUsed, ok := vlq.EncodeReversedTo(actualEncoded)
	if !ok {
		t.Error("Failed to encode LVLQ reversed")
		return
	}
	if bytesUsed != byteCount {
		t.Errorf("Expected to use %v bytes but actually used %v", byteCount, bytesUsed)
		return
	}
	if !bytes.Equal(expectedEncoded, actualEncoded) {
		t.Errorf("Expected reversed %v but got %v", expectedEncoded, actualEncoded)
		return
	}
	decoded, bytesUsed, ok := DecodeLvlqReversedFrom(actualEncoded)
	if !ok {
		t.Error("Failed to decode LVLQ")
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

func assertLvlqEncoded(t *testing.T, value int, expectedEncoded []byte) {
	vlq := Lvlq(value)
	byteCount := vlq.EncodedSize()
	actualEncoded := make([]byte, byteCount)
	bytesUsed, ok := vlq.EncodeTo(actualEncoded)
	if !ok {
		t.Error("Failed to encode LVLQ")
		return
	}
	if bytesUsed != byteCount {
		t.Errorf("Expected to use %v bytes but actually used %v", byteCount, bytesUsed)
		return
	}
	if !bytes.Equal(expectedEncoded, actualEncoded) {
		t.Errorf("Expected %v but got %v", expectedEncoded, actualEncoded)
	}
	decoded, bytesUsed, isComplete := DecodeLvlqFrom(actualEncoded)
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

	assertLvlqEncodedReversed(t, value, expectedEncoded)
}

func TestLvlq0(t *testing.T) {
	assertLvlqEncoded(t, 0, []byte{0})
}

func TestLvlq1(t *testing.T) {
	assertLvlqEncoded(t, 1, []byte{0xc0, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
}

func TestLvlq7F(t *testing.T) {
	assertLvlqEncoded(t, 0x7f, []byte{0xc0, 0xbf, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
}

func TestLvlq80(t *testing.T) {
	assertLvlqEncoded(t, 0x80, []byte{0xc0, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
}

func TestLvlqFF(t *testing.T) {
	assertLvlqEncoded(t, 0xff, []byte{0xc0, 0xff, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
}

func TestLvlq0x2005(t *testing.T) {
	assertLvlqEncoded(t, 0x2005, []byte{0xc0, 0x82, 0xa0, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
}

func TestLvlq0x3FFF(t *testing.T) {
	assertLvlqEncoded(t, 0x3fff, []byte{0xc0, 0xff, 0xbf, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
}

func TestLvlq0x4000(t *testing.T) {
	assertLvlqEncoded(t, 0x4000, []byte{0xc0, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
}

func TestLvlq0xFFFF(t *testing.T) {
	assertLvlqEncoded(t, 0xffff, []byte{0xc0, 0xff, 0xff, 0x81, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
}

func TestLvlq0x10000(t *testing.T) {
	assertLvlqEncoded(t, 0x10000, []byte{0x82, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
}

func TestLvlq0x1fffff(t *testing.T) {
	assertLvlqEncoded(t, 0x1fffff, []byte{0xc0, 0xff, 0xff, 0xbf, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})

}

func TestLvlq0x200000(t *testing.T) {
	assertLvlqEncoded(t, 0x200000, []byte{0xc0, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})

}
