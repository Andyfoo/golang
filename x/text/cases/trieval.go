// Code generated by running "go generate" in github.com/Andyfoo/golang/x/text. DO NOT EDIT.

package cases

// This file contains definitions for interpreting the trie value of the case
// trie generated by "go run gen*.go". It is shared by both the generator
// program and the resultant package. Sharing is achieved by the generator
// copying gen_trieval.go to trieval.go and changing what's above this comment.

// info holds case information for a single rune. It is the value returned
// by a trie lookup. Most mapping information can be stored in a single 16-bit
// value. If not, for example when a rune is mapped to multiple runes, the value
// stores some basic case data and an index into an array with additional data.
//
// The per-rune values have the following format:
//
//   if (exception) {
//     15..5  unsigned exception index
//         4  unused
//   } else {
//     15..8  XOR pattern or index to XOR pattern for case mapping
//            Only 13..8 are used for XOR patterns.
//         7  inverseFold (fold to upper, not to lower)
//         6  index: interpret the XOR pattern as an index
//            or isMid if case mode is cIgnorableUncased.
//      5..4  CCC: zero (normal or break), above or other
//   }
//      3  exception: interpret this value as an exception index
//         (TODO: is this bit necessary? Probably implied from case mode.)
//   2..0  case mode
//
// For the non-exceptional cases, a rune must be either uncased, lowercase or
// uppercase. If the rune is cased, the XOR pattern maps either a lowercase
// rune to uppercase or an uppercase rune to lowercase (applied to the 10
// least-significant bits of the rune).
//
// See the definitions below for a more detailed description of the various
// bits.
type info uint16

const (
	casedMask      = 0x0003
	fullCasedMask  = 0x0007
	ignorableMask  = 0x0006
	ignorableValue = 0x0004

	inverseFoldBit = 1 << 7
	isMidBit       = 1 << 6

	exceptionBit     = 1 << 3
	exceptionShift   = 5
	numExceptionBits = 11

	xorIndexBit = 1 << 6
	xorShift    = 8

	// There is no mapping if all xor bits and the exception bit are zero.
	hasMappingMask = 0xff80 | exceptionBit
)

// The case mode bits encodes the case type of a rune. This includes uncased,
// title, upper and lower case and case ignorable. (For a definition of these
// terms see Chapter 3 of The Unicode Standard Core Specification.) In some rare
// cases, a rune can be both cased and case-ignorable. This is encoded by
// cIgnorableCased. A rune of this type is always lower case. Some runes are
// cased while not having a mapping.
//
// A common pattern for scripts in the Unicode standard is for upper and lower
// case runes to alternate for increasing rune values (e.g. the accented Latin
// ranges starting from U+0100 and U+1E00 among others and some Cyrillic
// characters). We use this property by defining a cXORCase mode, where the case
// mode (always upper or lower case) is derived from the rune value. As the XOR
// pattern for case mappings is often identical for successive runes, using
// cXORCase can result in large series of identical trie values. This, in turn,
// allows us to better compress the trie blocks.
const (
	cUncased          info = iota // 000
	cTitle                        // 001
	cLower                        // 010
	cUpper                        // 011
	cIgnorableUncased             // 100
	cIgnorableCased               // 101 // lower case if mappings exist
	cXORCase                      // 11x // case is cLower | ((rune&1) ^ x)

	maxCaseMode = cUpper
)

func (c info) isCased() bool {
	return c&casedMask != 0
}

func (c info) isCaseIgnorable() bool {
	return c&ignorableMask == ignorableValue
}

func (c info) isNotCasedAndNotCaseIgnorable() bool {
	return c&fullCasedMask == 0
}

func (c info) isCaseIgnorableAndNotCased() bool {
	return c&fullCasedMask == cIgnorableUncased
}

func (c info) isMid() bool {
	return c&(fullCasedMask|isMidBit) == isMidBit|cIgnorableUncased
}

// The case mapping implementation will need to know about various Canonical
// Combining Class (CCC) values. We encode two of these in the trie value:
// cccZero (0) and cccAbove (230). If the value is cccOther, it means that
// CCC(r) > 0, but not 230. A value of cccBreak means that CCC(r) == 0 and that
// the rune also has the break category Break (see below).
const (
	cccBreak info = iota << 4
	cccZero
	cccAbove
	cccOther

	cccMask = cccBreak | cccZero | cccAbove | cccOther
)

const (
	starter       = 0
	above         = 230
	iotaSubscript = 240
)

// The exceptions slice holds data that does not fit in a normal info entry.
// The entry is pointed to by the exception index in an entry. It has the
// following format:
//
// Header
// byte 0:
//  7..6  unused
//  5..4  CCC type (same bits as entry)
//     3  unused
//  2..0  length of fold
//
// byte 1:
//   7..6  unused
//   5..3  length of 1st mapping of case type
//   2..0  length of 2nd mapping of case type
//
//   case     1st    2nd
//   lower -> upper, title
//   upper -> lower, title
//   title -> lower, upper
//
// Lengths with the value 0x7 indicate no value and implies no change.
// A length of 0 indicates a mapping to zero-length string.
//
// Body bytes:
//   case folding bytes
//   lowercase mapping bytes
//   uppercase mapping bytes
//   titlecase mapping bytes
//   closure mapping bytes (for NFKC_Casefold). (TODO)
//
// Fallbacks:
//   missing fold  -> lower
//   missing title -> upper
//   all missing   -> original rune
//
// exceptions starts with a dummy byte to enforce that there is no zero index
// value.
const (
	lengthMask = 0x07
	lengthBits = 3
	noChange   = 0
)

// References to generated trie.

var trie = newCaseTrie(0)

var sparse = sparseBlocks{
	values:  sparseValues[:],
	offsets: sparseOffsets[:],
}

// Sparse block lookup code.

// valueRange is an entry in a sparse block.
type valueRange struct {
	value  uint16
	lo, hi byte
}

type sparseBlocks struct {
	values  []valueRange
	offsets []uint16
}

// lookup returns the value from values block n for byte b using binary search.
func (s *sparseBlocks) lookup(n uint32, b byte) uint16 {
	lo := s.offsets[n]
	hi := s.offsets[n+1]
	for lo < hi {
		m := lo + (hi-lo)/2
		r := s.values[m]
		if r.lo <= b && b <= r.hi {
			return r.value
		}
		if b < r.lo {
			hi = m
		} else {
			lo = m + 1
		}
	}
	return 0
}

// lastRuneForTesting is the last rune used for testing. Everything after this
// is boring.
const lastRuneForTesting = rune(0x1FFFF)
