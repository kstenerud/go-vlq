// VLQ (Variable Length Quantity) is an unsigned integer encoding scheme
// designed for the MIDI file format.
//
// It encodes a value into a sequence of bytes where the lower 7 bits contain
// data and the high bit is used as a "continuation" bit. A decoder reads
// encoded bytes, filling a decoded unsigned integer 7 bits at a time in big
// endian order, until it encounters a byte with the high "continuation" bit
// cleared.
package vlq
