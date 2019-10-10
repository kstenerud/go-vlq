package vlq

// Copyright 2019 Karl Stenerud
// All rights reserved.
// Distributed under MIT license.

type Rvlq uint64

// Get the maximum value that can be stored in a RVLQ of the selected byte count
func MaxValueInBytes(byteCount int) uint64 {
	return (1 << (uint(byteCount) * 7)) - 1
}

// Get the number of bytes required to encode this RVLQ
func (this Rvlq) EncodedSize() int {
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

// Encode this RVLQ to a buffer.
// Returns false and the number of bytes it attempted to write if the buffer isn't big enough.
func (this Rvlq) EncodeTo(buffer []byte) (encodedByteCount int, ok bool) {
	value := this
	encodedByteCount = this.EncodedSize()
	if encodedByteCount > len(buffer) {
		return encodedByteCount, false
	}
	for i := encodedByteCount - 1; i >= 0; i-- {
		nextByte := byte(value >> uint(7*i) & 0x7f)
		if i > 0 {
			nextByte |= 0x80
		}
		buffer[encodedByteCount-i-1] = nextByte
	}
	return encodedByteCount, true
}

// Decode this RVLQ from a buffer. Returns true for isComplete once the VLQ
// is fully decoded (i.e. it has encountered a byte with the high bit cleared).
// This allows for progressive decoding of the VLQ value across multiple buffers.
func (this *Rvlq) DecodeFrom(buffer []byte) (decodedByteCount int, isComplete bool) {
	for _, v := range buffer {
		*this = *this<<7 | Rvlq(v&0x7f)
		decodedByteCount++
		if v&0x80 == 0 {
			return decodedByteCount, true
		}
	}
	return decodedByteCount, false
}

func DecodeRvlqFrom(buffer []byte) (value Rvlq, decodedByteCount int, isComplete bool) {
	decodedByteCount, isComplete = value.DecodeFrom(buffer)
	return value, decodedByteCount, isComplete
}

// Encode this RVLQ in reverse order to the end of a buffer.
// Returns false and the number of bytes it attempted to write if the buffer isn't big enough.
func (this Rvlq) EncodeReversedTo(buffer []byte) (encodedByteCount int, ok bool) {
	value := this
	encodedByteCount = this.EncodedSize()
	if encodedByteCount > len(buffer) {
		return encodedByteCount, false
	}
	start := len(buffer) - encodedByteCount
	for i := encodedByteCount - 1; i >= 0; i-- {
		nextByte := byte(value >> uint(7*i) & 0x7f)
		if i > 0 {
			nextByte |= 0x80
		}
		buffer[start+i] = nextByte
	}
	return encodedByteCount, true
}

// Decode this RVLQ in reverse order from the end of a buffer. Unlike DecodeFrom(),
// the reversed version must have all encoded bytes present in the buffer to
// decode.
func (this *Rvlq) DecodeReversedFrom(buffer []byte) (decodedByteCount int, ok bool) {
	for i := len(buffer) - 1; i >= 0; i-- {
		v := buffer[i]
		*this = *this<<7 | Rvlq(v&0x7f)
		decodedByteCount++
		if v&0x80 == 0 {
			return decodedByteCount, true
		}
	}
	return decodedByteCount, false
}

// Decode a RVLQ value in reverse order from the end of a buffer. Unlike DecodeFrom(),
// the reversed version must have all encoded bytes present in the buffer to
// decode.
func DecodeRvlqReversedFrom(buffer []byte) (value Rvlq, decodedByteCount int, ok bool) {
	decodedByteCount, ok = value.DecodeReversedFrom(buffer)
	return value, decodedByteCount, ok
}
