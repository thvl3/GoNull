package stealth

import (
	"fmt"
	"strings"
	"testing"
)

func TestEncodeToNullWidth(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "single character 'a'",
			input: "a", // ASCII 97, binary 01100001
			want:  ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthNonJoiner,
		},
		{
			name:  "string 'hi'",
			input: "hi", // h: 01101000, i: 01101001
			want: ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + // h
				ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthNonJoiner, // i
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeToNullWidth(tt.input); got != tt.want {
				t.Errorf("EncodeToNullWidth() = %v, want %v", got, tt.want)
				t.Errorf("EncodeToNullWidth() got (len %d): %s", len(got), got)
				t.Errorf("EncodeToNullWidth() want (len %d): %s", len(tt.want), tt.want)
			}
		})
	}
}

func TestDecodeFromNullWidth(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "decode 'a'",
			input:   ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthNonJoiner,
			want:    "a",
			wantErr: false,
		},
		{
			name: "decode 'hi'",
			input: ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + // h
				ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthNonJoiner, // i
			want:    "hi",
			wantErr: false,
		},
		{
			name:    "incomplete byte",
			input:   ZeroWidthSpace + ZeroWidthNonJoiner, // 01
			want:    "",
			wantErr: false, // Current implementation ignores incomplete bytes at the end
		},
		{
			name:    "mixed with other characters",
			input:   "A" + ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthNonJoiner + "B",
			want:    "a", // 'A' and 'B' should be ignored by current DecodeFromNullWidth
			wantErr: false,
		},
		{
			name: "contains other null-width chars like ZWJ, StartMarker, EndMarker",
			// The decoder should only pick up ZWS and ZWNJ
			input:   StartMarker + ZeroWidthSpace + ZeroWidthJoiner + ZeroWidthNonJoiner + EndMarker,
			want:    "", // Because ZWJ, Start, End are not ZWS/ZWNJ, effectively no bits for a byte
			wantErr: false,
		},
		{
			name:    "decode 'a' with markers and ZWJ (ignored)",
			input:   StartMarker + ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthNonJoiner + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthSpace + ZeroWidthNonJoiner + ZeroWidthJoiner + EndMarker,
			want:    "a",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeFromNullWidth(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeFromNullWidth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeFromNullWidth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty string", input: ""},
		{name: "simple string 'hello'", input: "hello"},
		{name: "string with spaces 'hello world'", input: "hello world"},
		{name: "string with numbers and symbols 'Cmd123!@#'", input: "Cmd123!@#"},
		{name: "longer string", input: "This is a longer test string to ensure encoding and decoding work correctly."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeToNullWidth(tt.input)
			decoded, err := DecodeFromNullWidth(encoded)

			if err != nil {
				t.Errorf("Roundtrip: DecodeFromNullWidth() returned error: %v for input '%s' (encoded: '%s')", err, tt.input, encoded)
				return
			}
			if decoded != tt.input {
				t.Errorf("Roundtrip: EncodeToNullWidth() then DecodeFromNullWidth()\nInput:    '%s'\nEncoded:  '%s' (len %d)\nDecoded:  '%s' (len %d)\nWant:     '%s' (len %d)", tt.input, encoded, len(encoded), decoded, len(decoded), tt.input, len(tt.input))
				// For easier debugging of character differences:
				if len(decoded) == len(tt.input) {
					for i := 0; i < len(tt.input); i++ {
						if tt.input[i] != decoded[i] {
							t.Errorf("Mismatch at index %d: input char %c (%d), decoded char %c (%d)", i, tt.input[i], tt.input[i], decoded[i], decoded[i])
						}
					}
				}
			}
		})
	}
}

func TestEmbedAndExtractText(t *testing.T) {
	command := "cmd"
	encodedCmd := EncodeToNullWidth(command)

	tests := []struct {
		name                    string
		originalText            string
		commandToEmbed          string // The raw command string
		expectedVisibleTextPart string // How the visible text should look after embedding (approx)
		wantExtractedCmd        string
		shouldFindData          bool // True if we expect ExtractFromText to find the markers and data
	}{
		{
			name:                    "simple embed: two words",
			originalText:            "Hello World",
			commandToEmbed:          command,
			expectedVisibleTextPart: "Hello " + StartMarker + encodedCmd + EndMarker + "World",
			wantExtractedCmd:        command,
			shouldFindData:          true,
		},
		{
			name:                    "embed: single word text",
			originalText:            "Hello",
			commandToEmbed:          command,
			expectedVisibleTextPart: StartMarker + encodedCmd + EndMarker + "Hello",
			wantExtractedCmd:        command,
			shouldFindData:          true,
		},
		{
			name:                    "embed: empty text",
			originalText:            "",
			commandToEmbed:          command,
			expectedVisibleTextPart: StartMarker + encodedCmd + EndMarker + "",
			wantExtractedCmd:        command,
			shouldFindData:          true,
		},
		{
			name:                    "embed: text with multiple spaces",
			originalText:            "First  Second Third", // Embed after "First "
			commandToEmbed:          command,
			expectedVisibleTextPart: "First " + StartMarker + encodedCmd + EndMarker + " Second Third",
			wantExtractedCmd:        command,
			shouldFindData:          true,
		},
		{
			name:                    "embed: empty command string",
			originalText:            "Some text",
			commandToEmbed:          "",          // encoded will be empty
			expectedVisibleTextPart: "Some text", // EmbedInText should return original text
			wantExtractedCmd:        "",
			shouldFindData:          false, // No markers embedded if command is empty
		},
		{
			name:                    "extract: no markers present",
			originalText:            "Just normal text without markers.",
			commandToEmbed:          "", // Not actually embedding, just setting up text for extraction test
			expectedVisibleTextPart: "Just normal text without markers.",
			wantExtractedCmd:        "",
			shouldFindData:          false,
		},
		{
			name:                    "extract: only start marker",
			originalText:            "Text with " + StartMarker + " unmatched.",
			commandToEmbed:          "",
			expectedVisibleTextPart: "Text with " + StartMarker + " unmatched.",
			wantExtractedCmd:        "",
			shouldFindData:          false,
		},
		{
			name:                    "extract: only end marker",
			originalText:            "Text with " + EndMarker + " unmatched.",
			commandToEmbed:          "",
			expectedVisibleTextPart: "Text with " + EndMarker + " unmatched.",
			wantExtractedCmd:        "",
			shouldFindData:          false,
		},
		{
			name:                    "extract: markers in wrong order",
			originalText:            "Text with " + EndMarker + encodedCmd + StartMarker + " wrong order.",
			commandToEmbed:          command, // command used to generate encodedCmd for the text
			expectedVisibleTextPart: "Text with " + EndMarker + encodedCmd + StartMarker + " wrong order.",
			wantExtractedCmd:        "",
			shouldFindData:          false,
		},
		{
			name: "extract: data with other null-width chars between markers (should be filtered)",
			// Construct text manually for this test case
			originalText:            "Before" + StartMarker + ZeroWidthSpace + ZeroWidthJoiner + ZeroWidthNonJoiner + EndMarker + "After",
			commandToEmbed:          "", // We are not using EmbedInText for this specific case
			expectedVisibleTextPart: "Before" + StartMarker + ZeroWidthSpace + ZeroWidthJoiner + ZeroWidthNonJoiner + EndMarker + "After",
			// Expected extracted bits: ZWS ZWNJ (01). This is not enough for a full byte.
			// DecodeFromNullWidth will return "" for incomplete bytes.
			// The raw extracted null width sequence by ExtractFromText should be ZWS + ZWNJ.
			wantExtractedCmd: "",
			shouldFindData:   true, // Markers are found, payload is ZWS+ZWNJ
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var textToProcess string
			encodedCommandForTest := EncodeToNullWidth(tt.commandToEmbed)

			if tt.commandToEmbed != "" || (tt.name != "extract: data with other null-width chars between markers (should be filtered)" && strings.Contains(tt.name, "embed:")) {
				// Use EmbedInText for embedding tests or if commandToEmbed is meant to be embedded
				textToProcess = EmbedInText(tt.originalText, encodedCommandForTest)
				// Verify visible part for embedding tests
				if strings.Contains(tt.name, "embed:") && textToProcess != tt.expectedVisibleTextPart {
					t.Errorf("EmbedInText() output mismatch\nGot:  '%s'\nWant: '%s'", textToProcess, tt.expectedVisibleTextPart)
				}
			} else {
				// For specific extraction tests where text is manually constructed
				textToProcess = tt.originalText
			}

			extractedNullWidthData := ExtractFromText(textToProcess)

			if tt.shouldFindData {
				// If we expect to find data, extractedNullWidthData should be the encoded command
				// or the specific sequence like ZWS+ZWNJ in the mixed chars test
				expectedRawExtraction := encodedCommandForTest
				if tt.name == "extract: data with other null-width chars between markers (should be filtered)" {
					expectedRawExtraction = ZeroWidthSpace + ZeroWidthNonJoiner // ZWJ is filtered out by ExtractFromText
				}

				if extractedNullWidthData == "" && tt.commandToEmbed != "" { // Allow empty if commandToEmbed was empty
					// only fail if we expected a non-empty command but got empty extraction
					t.Errorf("ExtractFromText() returned empty string, expected to find data for command '%s' from text '%s'", tt.commandToEmbed, textToProcess)
					return
				}
				if extractedNullWidthData != expectedRawExtraction {
					t.Errorf("ExtractFromText() raw data mismatch\nGot:  '%s' (len %d)\nWant: '%s' (len %d)\nFrom Text: '%s'",
						extractedNullWidthData, len(extractedNullWidthData), expectedRawExtraction, len(expectedRawExtraction), textToProcess)
				}

				decodedCommand, err := DecodeFromNullWidth(extractedNullWidthData)
				if err != nil {
					t.Errorf("DecodeFromNullWidth() after extract failed with error: %v for extracted '%s'", err, extractedNullWidthData)
					return
				}
				if decodedCommand != tt.wantExtractedCmd {
					t.Errorf("Decoded command after extract does not match.\nGot:  '%s'\nWant: '%s'\nExtracted Raw: '%s'", decodedCommand, tt.wantExtractedCmd, extractedNullWidthData)
				}

			} else { // shouldNotFindData
				if extractedNullWidthData != "" {
					t.Errorf("ExtractFromText() expected empty string but got '%s' from text '%s'", extractedNullWidthData, textToProcess)
				}
				// Also ensure decoding an empty string (or whatever was extracted) results in the wanted empty command
				decodedCommand, err := DecodeFromNullWidth(extractedNullWidthData)
				if err != nil {
					t.Errorf("DecodeFromNullWidth() on expected empty extraction failed: %v", err)
				}
				if decodedCommand != tt.wantExtractedCmd { // Should be ""
					t.Errorf("Decoded command from expected non-data scenario mismatch. Got '%s', want '%s'", decodedCommand, tt.wantExtractedCmd)
				}
			}
		})
	}
}

// Helper to visualize binary for debugging if needed
func printBinary(s string) {
	for _, c := range s {
		fmt.Printf("%08b ", c)
	}
	fmt.Println()
}
