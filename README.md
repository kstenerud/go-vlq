Variable Length Quantity
========================

A go implementation of the [VLQ type](https://github.com/kstenerud/vlq/blob/master/vlq-specification.md). VLQ is an encoding scheme to compress unsigned integers.


Library Usage
-------------

```golang
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
```

Result:

    Encoded 30000 into the byte sequence [129 234 48]
    Decoded value 30000 from 3 bytes


```golang
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
```

Result:

    Used 2 bytes decoding VLQ value. Is complete = false
    Used 1 bytes decoding VLQ value. Is complete = true
    Decoded value: 30000


```golang
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
```

Result:

    Encoded 30000 into the byte sequence [0 0 0 0 0 0 0 48 234 129]
    Decoded value 30000 from 3 bytes
