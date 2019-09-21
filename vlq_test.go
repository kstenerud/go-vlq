package vlq

import (
	"bytes"
	"testing"
)

func TestExtend(t *testing.T) {
	encodeCount := 2
	expected := []byte{0x80, 0x80, 0, 0, 0}
	actual := make([]byte, 5)
	bytesEncoded, err := Extend(actual, encodeCount)
	if err != nil {
		t.Errorf("%v", err)
	}
	if bytesEncoded != encodeCount {
		t.Errorf("Expected to encode %v bytes but got %v", encodeCount, bytesEncoded)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected buffer %v but got %v", expected, actual)
	}
}

func TestIsExtended(t *testing.T) {
	buffer := []byte{0x80, 0x00}
	isExtended := IsExtended(buffer)
	if !isExtended {
		t.Errorf("Expected extended but was not extended")
	}
}

func TestIsNotExtended(t *testing.T) {
	buffer := []byte{0x00, 0x00}
	isExtended := IsExtended(buffer)
	if isExtended {
		t.Errorf("Expected not extended but was extended")
	}
}
