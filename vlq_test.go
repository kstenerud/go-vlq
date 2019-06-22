package vlq

import (
	"bytes"
	"fmt"
	"testing"
)

func assertEncodedReversed(t *testing.T, value int, expectedNonReversedEncoded []byte) {
	byteCount := len(expectedNonReversedEncoded)
	inflatedByteCount := byteCount + 10
	expectedEncoded := make([]byte, inflatedByteCount)
	for i, v := range expectedNonReversedEncoded {
		expectedEncoded[inflatedByteCount-i-1] = v
	}
	actualEncoded := make([]byte, inflatedByteCount)
	vlq := Vlq(value)

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
	decoded, bytesUsed, err := DecodeReversedFrom(actualEncoded)
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

func assertEncoded(t *testing.T, value int, expectedEncoded []byte) {
	vlq := Vlq(value)
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
	decoded, bytesUsed, isComplete := DecodeFrom(actualEncoded)
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

	assertEncodedReversed(t, value, expectedEncoded)
}

func Test0(t *testing.T) {
	assertEncoded(t, 0, []byte{0})
}

func Test1(t *testing.T) {
	assertEncoded(t, 1, []byte{1})
}

func Test7F(t *testing.T) {
	assertEncoded(t, 0x7f, []byte{0x7f})
}

func Test80(t *testing.T) {
	assertEncoded(t, 0x80, []byte{0x81, 0x00})
}

func TestFF(t *testing.T) {
	assertEncoded(t, 0xff, []byte{0x81, 0x7f})
}

func Test2005(t *testing.T) {
	assertEncoded(t, 0x2005, []byte{0xc0, 0x05})
}

func Test3FFF(t *testing.T) {
	assertEncoded(t, 0x3fff, []byte{0xff, 0x7f})
}

func Test4000(t *testing.T) {
	assertEncoded(t, 0x4000, []byte{0x81, 0x80, 0x00})
}

func TestFFFF(t *testing.T) {
	assertEncoded(t, 0xffff, []byte{0x83, 0xff, 0x7f})
}

func Test10000(t *testing.T) {
	assertEncoded(t, 0x10000, []byte{0x84, 0x80, 0x00})
}

func Test1fffff(t *testing.T) {
	assertEncoded(t, 0x1fffff, []byte{0xff, 0xff, 0x7f})

}

func Test200000(t *testing.T) {
	assertEncoded(t, 0x200000, []byte{0x81, 0x80, 0x80, 0x00})

}

func readme_example_normal() {
	// Standard operation

	value := Vlq(30000)
	byteCount := value.EncodedSize()
	buffer := make([]byte, byteCount)
	value.EncodeTo(buffer)
	fmt.Printf("Encoded %v into the byte sequence %v\n", value, buffer)

	decodedValue, bytesUsed, isComplete := DecodeFrom(buffer)
	if !isComplete {
		// TODO: Error handling
	}
	fmt.Printf("Decoded value %v from %v bytes\n", decodedValue, bytesUsed)
}

func readme_example_split_reads() {
	// Simulate the VLQ being split across two reads

	value := Vlq(30000)
	byteCount := value.EncodedSize()
	buffer := make([]byte, byteCount)
	value.EncodeTo(buffer)

	buffer1 := buffer[:2]
	buffer2 := buffer[2:]

	var decodedValue Vlq
	bytesUsed, isComplete := decodedValue.DecodeFrom(buffer1)
	fmt.Printf("Used %v bytes decoding VLQ value. Is complete = %v\n", bytesUsed, isComplete)

	bytesUsed, isComplete = decodedValue.DecodeFrom(buffer2)
	fmt.Printf("Used %v bytes decoding VLQ value. Is complete = %v\n", bytesUsed, isComplete)

	fmt.Printf("Decoded value: %v\n", decodedValue)
}

func readme_example_reversed() {
	// Reversed mode, which allows reading the value from the end of a buffer

	value := Vlq(30000)
	buffer := make([]byte, 10)
	value.EncodeReversedTo(buffer)
	fmt.Printf("Encoded %v into the byte sequence %v\n", value, buffer)

	decodedValue, bytesUsed, err := DecodeReversedFrom(buffer)
	if err != nil {
		// TODO: Error handling
	}
	fmt.Printf("Decoded value %v from %v bytes\n", decodedValue, bytesUsed)
}

func TestExamples(t *testing.T) {
	readme_example_normal()
	readme_example_split_reads()
	readme_example_reversed()
}
