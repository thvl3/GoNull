package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gonull/internal/common/stealth"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// 1. Get the text containing the hidden command
	fmt.Println("Paste the text containing the hidden command (e.g., copied from a webpage or server output):")
	// Read multiple lines until a specific terminator or EOF? For now, assume single line paste or handle newline.
	// A simple ReadString might be problematic if the text has newlines naturally.
	// Let's read line by line until an empty line for simplicity in pasting multi-line text.
	fmt.Println("(End input with an empty line or Ctrl+D)")

	var inputTextBuilder strings.Builder
	for {
		line, err := reader.ReadString('\n')
		trimmedLine := strings.TrimSpace(line)

		if err != nil {
			if err.Error() == "EOF" { // Handle Ctrl+D
				break
			}
			log.Fatalf("Failed to read input text: %v", err)
		}

		if inputTextBuilder.Len() > 0 {
			inputTextBuilder.WriteString("\n") // Add back newline if it's not the first line of pasted content
		}
		inputTextBuilder.WriteString(strings.TrimRight(line, "\r\n")) // Add the line without its own newline

		// For simple single-line paste, this check might be too aggressive.
		// For now, let's assume if a trimmed line is empty, it's the end of multi-line paste.
		if trimmedLine == "" && inputTextBuilder.Len() > 0 { // Check if it was truly an empty line signal
			break
		}
		// If the first line entered is empty, also break.
		if trimmedLine == "" && inputTextBuilder.Len() == 0 {
			break
		}
	}

	inputText := strings.TrimSpace(inputTextBuilder.String())

	if inputText == "" {
		log.Println("No input text provided. Exiting.")
		return
	}

	// 2. Extract the null-width characters (the encoded command bits)
	extractedBits := stealth.ExtractFromText(inputText)

	if extractedBits == "" {
		log.Println("No hidden command found in the provided text (no markers or empty payload).")
		log.Printf("Input text was: %q", inputText) // Log the input for debugging
		return
	}

	// 3. Decode the command
	decodedCommand, err := stealth.DecodeFromNullWidth(extractedBits)
	if err != nil {
		log.Fatalf("Failed to decode command from extracted bits: %v\nExtracted bits: %s", err, extractedBits)
	}

	if decodedCommand == "" {
		log.Println("Decoded command is empty. This might mean the extracted bits were insufficient or invalid.")
		log.Printf("Extracted bits were: %s", extractedBits)
		// Potentially, this is a valid outcome if an empty command was sent.
		// However, our server prevents encoding truly empty commands, though the payload bits might be non-byte-aligned.
	}

	// 4. Print the decoded command
	fmt.Println("\n--- Decoded Command ---")
	fmt.Printf("Command: '%s'\n", decodedCommand)
	fmt.Println("--- End of Command --- ")

	// TODO: Add command execution logic here
	log.Println("\nNext step: Implement command execution for the decoded command.")
}
