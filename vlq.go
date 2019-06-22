// VLQ (Variable Length Quantity) is an unsigned integer encoding scheme
// designed for the MIDI file format.
//
// It encodes a value into a sequence of bytes where the lower 7 bits contain
// data and the high bit is used as a "continuation" bit. A decoder reads
// encoded bytes, filling a decoded unsigned integer 7 bits at a time in big
// endian order, until it encounters a byte with the high "continuation" bit
// cleared.
package vlq

// Copyright 2019 Karl Stenerud
// All rights reserved.
// Distributed under MIT license.

import (
	"fmt"
)

type Vlq uint64

// Get the maximum value that can be stored in a VLQ of the selected byte count
func MaxValueInBytes(byteCount int) uint64 {
	return (1 << (uint(byteCount) * 7)) - 1
}

// Get the number of bytes required to encode this VLQ
func (this Vlq) EncodedSize() int {
	value := this
	if value <= 0x7f {
		return 1
	}
	size := 0
	for value > 0 {
		value >>= 7
		size++
	}
	return size
}

// Encode this VLQ to a buffer. Returns an error if the buffer isn't big enough.
func (this Vlq) EncodeTo(buffer []byte) (bytesEncoded int, err error) {
	value := this
	bytesEncoded = this.EncodedSize()
	if bytesEncoded > len(buffer) {
		return 0, fmt.Errorf("%v free bytes required to encode, but only %v available", this.EncodedSize(), len(buffer))
	}
	for i := bytesEncoded - 1; i >= 0; i-- {
		nextByte := byte(value >> uint(7*i) & 0x7f)
		if i > 0 {
			nextByte |= 0x80
		}
		buffer[bytesEncoded-i-1] = nextByte
	}
	return bytesEncoded, nil
}

// Decode this VLQ from a buffer. Returns true for isComplete once the VLQ
// is fully decoded (i.e. it has encountered a byte with the high bit cleared).
// This allows for progressive decoding of the VLQ value across multiple buffers.
func (this *Vlq) DecodeFrom(buffer []byte) (bytesDecoded int, isComplete bool) {
	for _, v := range buffer {
		*this = *this<<7 | Vlq(v&0x7f)
		bytesDecoded++
		if v&0x80 == 0 {
			return bytesDecoded, true
		}
	}
	return bytesDecoded, false
}

func DecodeFrom(buffer []byte) (value Vlq, bytesDecoded int, isComplete bool) {
	bytesDecoded, isComplete = value.DecodeFrom(buffer)
	return value, bytesDecoded, isComplete
}

// Encode this VLQ in reverse order to the end of a buffer.
// Returns an error if the buffer isn't big enough.
func (this Vlq) EncodeReversedTo(buffer []byte) (bytesEncoded int, err error) {
	value := this
	bytesEncoded = this.EncodedSize()
	if bytesEncoded > len(buffer) {
		return 0, fmt.Errorf("%v free bytes required to encode, but only %v available", this.EncodedSize(), len(buffer))
	}
	start := len(buffer) - bytesEncoded
	for i := bytesEncoded - 1; i >= 0; i-- {
		nextByte := byte(value >> uint(7*i) & 0x7f)
		if i > 0 {
			nextByte |= 0x80
		}
		buffer[start+i] = nextByte
	}
	return bytesEncoded, nil
}

// Decode this VLQ in reverse order from the end of a buffer. Unlike DecodeFrom(),
// the reversed version must have all encoded bytes present in the buffer to
// decode.
func (this *Vlq) DecodeReversedFrom(buffer []byte) (bytesDecoded int, err error) {
	for i := len(buffer) - 1; i >= 0; i-- {
		v := buffer[i]
		*this = *this<<7 | Vlq(v&0x7f)
		bytesDecoded++
		if v&0x80 == 0 {
			return bytesDecoded, nil
		}
	}
	return 0, fmt.Errorf("Buffer does not contain the complete encoded reverse VLQ value")
}

// Decode a VLQ value in reverse order from the end of a buffer. Unlike DecodeFrom(),
// the reversed version must have all encoded bytes present in the buffer to
// decode.
func DecodeReversedFrom(buffer []byte) (value Vlq, bytesDecoded int, err error) {
	bytesDecoded, err = value.DecodeReversedFrom(buffer)
	return value, bytesDecoded, err
}
