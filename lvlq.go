package vlq

// Copyright 2019 Karl Stenerud
// All rights reserved.
// Distributed under MIT license.

type Lvlq uint64

// Get the number of bytes required to encode this RVLQ
func (this Lvlq) EncodedSize() int {
	value := this
	if value == 0 {
		return 1
	}
	size := 0
	for value > 0 {
		value <<= 7
		size++
	}
	return size
}

// Encode this LVLQ to a buffer.
// Returns false and the number of bytes it attempted to write if the buffer isn't big enough.
func (this Lvlq) EncodeTo(buffer []byte) (encodedByteCount int, ok bool) {
	value := this
	encodedByteCount = this.EncodedSize()
	if encodedByteCount > len(buffer) {
		return encodedByteCount, false
	}
	if value == 0 {
		buffer[0] = 0
		return 1, true
	}

	groupCount := encodedByteCount
	index := 0

	extraBitCount := uint(64) % 7
	extraMask := Lvlq(1<<extraBitCount) - 1

	if (value & extraMask) != 0 {
		extraShift := (7 - extraBitCount) % 7
		buffer[index] = (uint8(value&extraMask) << extraShift) | 0x80
		value >>= extraBitCount
		index++
	} else {
		value >>= uint(64 - groupCount*7)
	}

	for ; index < groupCount; index++ {
		nextByte := uint8(value & 0x7f)
		if index < groupCount-1 {
			nextByte |= 0x80
		}
		buffer[index] = nextByte
		value >>= 7
	}
	return encodedByteCount, true
}

// Decode this LVLQ from a buffer. Returns true for isComplete once the VLQ
// is fully decoded (i.e. it has encountered a byte with the high bit cleared).
// This allows for progressive decoding of the VLQ value across multiple buffers.
func (this *Lvlq) DecodeFrom(buffer []byte) (decodedByteCount int, isComplete bool) {
	for _, v := range buffer {
		*this = *this>>7 | (Lvlq(v&0x7f) << 57)
		decodedByteCount++
		if v&0x80 == 0 {
			return decodedByteCount, true
		}
	}
	return decodedByteCount, false
}

func DecodeLvlqFrom(buffer []byte) (value Lvlq, decodedByteCount int, isComplete bool) {
	decodedByteCount, isComplete = value.DecodeFrom(buffer)
	return value, decodedByteCount, isComplete
}

// Encode this LVLQ in reverse order to the end of a buffer.
// Returns false and the number of bytes it attempted to write if the buffer isn't big enough.
func (this Lvlq) EncodeReversedTo(buffer []byte) (encodedByteCount int, ok bool) {
	value := this
	encodedByteCount = this.EncodedSize()
	if encodedByteCount > len(buffer) {
		return encodedByteCount, false
	}
	start := len(buffer) - encodedByteCount
	if value == 0 {
		buffer[start] = 0
		return 1, true
	}

	groupCount := encodedByteCount
	index := len(buffer) - 1

	extraBitCount := uint(64) % 7
	extraMask := Lvlq(1<<extraBitCount) - 1

	if (value & extraMask) != 0 {
		extraShift := (7 - extraBitCount) % 7
		buffer[index] = (uint8(value&extraMask) << extraShift) | 0x80
		value >>= extraBitCount
		index--
	} else {
		value >>= uint(64 - groupCount*7)
	}

	for ; index > start; index-- {
		nextByte := uint8(value & 0x7f)
		if index > 0 {
			nextByte |= 0x80
		}
		buffer[index] = nextByte
		value >>= 7
	}

	return encodedByteCount, true
}

// Decode this RVLQ in reverse order from the end of a buffer. Unlike DecodeFrom(),
// the reversed version must have all encoded bytes present in the buffer to
// decode.
func (this *Lvlq) DecodeReversedFrom(buffer []byte) (decodedByteCount int, ok bool) {
	for i := len(buffer) - 1; i >= 0; i-- {
		v := buffer[i]
		*this = *this>>7 | (Lvlq(v&0x7f) << 57)
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
func DecodeLvlqReversedFrom(buffer []byte) (value Lvlq, decodedByteCount int, ok bool) {
	decodedByteCount, ok = value.DecodeReversedFrom(buffer)
	return value, decodedByteCount, ok
}
