package vlq

import (
	"fmt"
	"testing"
)

func readme_example_normal() {
	// Standard operation

	value := Rvlq(30000)
	byteCount := value.EncodedSize()
	buffer := make([]byte, byteCount)
	value.EncodeTo(buffer)
	fmt.Printf("Encoded %v into the byte sequence %v\n", value, buffer)

	decodedValue, bytesUsed, isComplete := DecodeRvlqFrom(buffer)
	if !isComplete {
		// TODO: Error handling
	}
	fmt.Printf("Decoded value %v from %v bytes\n", decodedValue, bytesUsed)
}

func readme_example_split_reads() {
	// Simulate the VLQ being split across two reads

	value := Rvlq(30000)
	byteCount := value.EncodedSize()
	buffer := make([]byte, byteCount)
	value.EncodeTo(buffer)

	buffer1 := buffer[:2]
	buffer2 := buffer[2:]

	var decodedValue Rvlq
	bytesUsed, isComplete := decodedValue.DecodeFrom(buffer1)
	fmt.Printf("Used %v bytes decoding VLQ value. Is complete = %v\n", bytesUsed, isComplete)

	bytesUsed, isComplete = decodedValue.DecodeFrom(buffer2)
	fmt.Printf("Used %v bytes decoding VLQ value. Is complete = %v\n", bytesUsed, isComplete)

	fmt.Printf("Decoded value: %v\n", decodedValue)
}

func readme_example_reversed() {
	// Reversed mode, which allows reading the value from the end of a buffer

	value := Rvlq(30000)
	buffer := make([]byte, 10)
	value.EncodeReversedTo(buffer)
	fmt.Printf("Encoded %v into the byte sequence %v\n", value, buffer)

	decodedValue, bytesUsed, ok := DecodeRvlqReversedFrom(buffer)
	if !ok {
		// TODO: Didn't find end of sequence
	}
	fmt.Printf("Decoded value %v from %v bytes\n", decodedValue, bytesUsed)
}

func TestExamples(t *testing.T) {
	readme_example_normal()
	readme_example_split_reads()
	readme_example_reversed()
}
