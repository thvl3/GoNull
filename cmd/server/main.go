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

	// 1. Get command from user
	fmt.Print("Enter command to encode: ")
	command, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read command: %v", err)
	}
	command = strings.TrimSpace(command)

	if command == "" {
		log.Println("No command entered. Exiting.")
		return
	}

	// 2. Get cover text from user
	fmt.Print("Enter cover text (e.g., a sentence or paragraph): ")
	coverText, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read cover text: %v", err)
	}
	coverText = strings.TrimSpace(coverText)

	if coverText == "" {
		log.Println("No cover text entered. Using a default placeholder.")
		coverText = "This is some default cover text."
	}

	// 3. Encode the command
	nullWidthEncodedCommand := stealth.EncodeToNullWidth(command)
	if nullWidthEncodedCommand == "" && command != "" {
		log.Printf("Warning: Encoding the command '%s' resulted in an empty null-width string. This is unusual.", command)
	}

	// 4. Embed into cover text
	textWithHiddenCommand := stealth.EmbedInText(coverText, nullWidthEncodedCommand)

	// 5. Print the result
	fmt.Println("\n--- Text with Embedded Command ---")
	fmt.Println(textWithHiddenCommand)
	fmt.Println("--- End of Text --- reproducing it exactly for verification ---")
	// Print with quotes to clearly show leading/trailing spaces if any, and the exact string
	// This helps in verifying that what's printed is what you can copy
	fmt.Printf("%q\n", textWithHiddenCommand)

	log.Println("\nInstructions: Copy the text between '--- Text with Embedded Command ---' and '--- End of Text ---' (or the quoted version below it). This text can now be placed on a webpage or social media post. The client will attempt to extract the command from it.")

	// For demonstration, let's also try to extract and decode it here to verify
	extractedBits := stealth.ExtractFromText(textWithHiddenCommand)
	decodedCommand, err := stealth.DecodeFromNullWidth(extractedBits)
	if err != nil {
		log.Printf("\nServer-side verification: Error decoding extracted bits: %v", err)
		log.Printf("Extracted bits were: %s", extractedBits)
	} else {
		if decodedCommand == command {
			log.Printf("\nServer-side verification: Successfully decoded command: '%s'", decodedCommand)
		} else {
			log.Printf("\nServer-side verification: Mismatch! Original: '%s', Decoded: '%s'", command, decodedCommand)
			log.Printf("Original encoded: %s", nullWidthEncodedCommand)
			log.Printf("Extracted bits:   %s", extractedBits)
		}
	}
}
