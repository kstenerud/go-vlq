// VLQ (Variable Length Quantity) is an unsigned integer encoding scheme
// designed for the MIDI file format.
//
// It encodes a value into a sequence of bytes where the lower 7 bits contain
// data and the high bit is used as a "continuation" bit. A decoder reads
// encoded bytes, filling a decoded unsigned integer 7 bits at a time in big
// endian order, until it encounters a byte with the high "continuation" bit
// cleared.
package vlq

import (
	"fmt"
)

// Maximum byte length that this library will encode (for a 64-bit value)
const MaxEncodeLength = 10

// Returns true if the encoded VLQ is extended (prepended with 0x80)
func IsExtended(buffer []byte) bool {
	return buffer[0] == 0x80
}

// Adds VLQ extension groups to a buffer.
func Extend(buffer []byte, groupCount int) (bytesEncoded int, err error) {
	if len(buffer) < groupCount {
		return 0, fmt.Errorf("Require %d bytes but only %d available", groupCount, len(buffer))
	}
	for i := 0; i < groupCount; i++ {
		buffer[i] = 0x80
	}
	return groupCount, nil
}
