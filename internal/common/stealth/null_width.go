package stealth

import (
	"errors"
	"strings"
)

// Define the null-width characters we'll use
const (
	ZeroWidthSpace     = "\u200b" // Represents 0
	ZeroWidthNonJoiner = "\u200c" // Represents 1
	ZeroWidthJoiner    = "\u200d" // Could be used as a separator or for other purposes
	StartMarker        = "\u200e" // LeftToRightMark (LRM)
	EndMarker          = "\u200f" // RightToLeftMark (RLM)
)

// EncodeToNullWidth takes a string and encodes it into a sequence of null-width characters.
// Each byte of the input string is converted to its 8-bit binary representation.
// ZeroWidthSpace represents '0' and ZeroWidthNonJoiner represents '1'.
func EncodeToNullWidth(data string) string {
	var encoded strings.Builder

	for _, b := range []byte(data) {
		// Convert byte to its 8-bit binary representation
		// Example: 'a' (97) -> "01100001"
		binary := byteToBinary(b)
		for _, bit := range binary {
			if bit == '0' {
				encoded.WriteString(ZeroWidthSpace)
			} else {
				encoded.WriteString(ZeroWidthNonJoiner)
			}
		}
	}
	return encoded.String()
}

// DecodeFromNullWidth takes a string containing null-width characters and decodes it.
// It reconstructs the original string from the binary representation.
func DecodeFromNullWidth(encodedData string) (string, error) {
	var binaryString strings.Builder
	var decodedBytes []byte

	for _, r := range encodedData {
		switch string(r) {
		case ZeroWidthSpace:
			binaryString.WriteRune('0')
		case ZeroWidthNonJoiner:
			binaryString.WriteRune('1')
		case ZeroWidthJoiner, StartMarker, EndMarker:
			// Ignore other null-width characters for now, or handle as delimiters/markers later
			continue
		default:
			// If we encounter a non-null-width character, it might mean the end of our data
			// or an error. For now, we'll assume it's not part of our encoded message.
			// This behavior might need to be refined based on how text is fetched from pages.
			continue
		}

		// If we have 8 bits, convert to byte
		if binaryString.Len() == 8 {
			b, err := binaryToByte(binaryString.String())
			if err != nil {
				// This should ideally not happen if our encoding is correct
				// and we only feed it ZeroWidthSpace/NonJoiner
				return "", err
			}
			decodedBytes = append(decodedBytes, b)
			binaryString.Reset()
		}
	}

	// Handle any remaining bits (should be less than 8)
	// This indicates an incomplete byte, which might be an error or requires padding logic
	if binaryString.Len() > 0 {
		// For now, we'll consider this an error or incomplete message.
		// Depending on the protocol, we might need padding or a specific end marker.
		// log.Printf("Warning: Incomplete byte at the end of decoding: %s", binaryString.String())
	}

	return string(decodedBytes), nil
}

// byteToBinary converts a byte to its 8-bit string representation.
// e.g., 97 ('a') -> "01100001"
func byteToBinary(b byte) string {
	var out strings.Builder
	for i := 7; i >= 0; i-- {
		if (b>>i)&1 == 1 {
			out.WriteRune('1')
		} else {
			out.WriteRune('0')
		}
	}
	return out.String()
}

// binaryToByte converts an 8-bit string representation to a byte.
// e.g., "01100001" -> 97 ('a')
func binaryToByte(binStr string) (byte, error) {
	if len(binStr) != 8 {
		return 0, errors.New("binary string must be 8 bits long")
	}
	var val byte
	for i := 0; i < 8; i++ {
		if binStr[i] == '1' {
			val |= (1 << (7 - i))
		}
	}
	return val, nil
}

// EmbedInText takes regular text and embeds the null-width encoded data within it.
// This is a refined embedding strategy using start and end markers.
// It inserts the payload after the first word of the text, or at the beginning if no spaces.
func EmbedInText(text, encodedData string) string {
	if encodedData == "" {
		return text // Nothing to embed
	}
	payload := StartMarker + encodedData + EndMarker

	firstSpaceIndex := strings.Index(text, " ")
	if firstSpaceIndex != -1 {
		// Insert after the first word
		return text[:firstSpaceIndex+1] + payload + text[firstSpaceIndex+1:]
	} else {
		// No spaces, or empty text; prepend/append based on preference.
		// Prepending might be less conspicuous if text is short or a single identifier.
		// For more general cases, one might append or intersperse more deeply.
		// Let's try inserting at the beginning for simplicity here.
		return payload + text
	}
}

// ExtractFromText attempts to find and extract the core payload (ZWS and ZWNJ)
// from a larger body of text, delineated by StartMarker and EndMarker.
func ExtractFromText(textWithHiddenData string) string {
	startIndex := strings.Index(textWithHiddenData, StartMarker)
	if startIndex == -1 {
		return "" // Start marker not found
	}

	// Search for EndMarker *after* the StartMarker
	searchArea := textWithHiddenData[startIndex+len(StartMarker):]
	endIndexInSearchArea := strings.Index(searchArea, EndMarker)

	if endIndexInSearchArea == -1 {
		return "" // End marker not found after start marker
	}

	// The raw data is between the end of StartMarker and the beginning of EndMarker
	rawData := searchArea[:endIndexInSearchArea]

	// Filter rawData to keep only ZeroWidthSpace and ZeroWidthNonJoiner
	var extractedPayload strings.Builder
	for _, r := range rawData {
		switch string(r) {
		case ZeroWidthSpace, ZeroWidthNonJoiner:
			extractedPayload.WriteRune(r)
			// Potentially log or handle other unexpected null-width chars if desired
			// case ZeroWidthJoiner, StartMarker, EndMarker:
			//   log.Printf("Warning: Unexpected marker or ZWJ inside payload: %q", r)
		}
	}

	return extractedPayload.String()
}
