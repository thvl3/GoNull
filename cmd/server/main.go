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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // More informative logging
	log.Println("Server starting...")

	reader := bufio.NewReader(os.Stdin)

	// 1. Get command from user
	fmt.Print("Enter command to encode: ")
	command, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading command from stdin: %v", err)
	}
	command = strings.TrimSpace(command)

	if command == "" {
		log.Println("No command entered by user. Exiting.")
		return
	}
	log.Printf("Command to encode: '%s'", command)

	// 2. Get cover text from user
	fmt.Print("Enter cover text (e.g., a sentence or paragraph): ")
	coverText, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading cover text from stdin: %v", err)
	}
	coverText = strings.TrimSpace(coverText)

	if coverText == "" {
		log.Println("No cover text entered. Using a default placeholder text.")
		coverText = "This is some default cover text."
	}
	log.Printf("Using cover text (first 50 chars if long): '%.50s...'", coverText)

	// 3. Encode the command
	log.Printf("Encoding command '%s' into null-width characters...", command)
	nullWidthEncodedCommand := stealth.EncodeToNullWidth(command)
	if nullWidthEncodedCommand == "" && command != "" {
		// This case should ideally not happen with current encoding which processes each byte.
		log.Printf("Warning: Encoding the command '%s' resulted in an empty null-width string. This is unexpected.", command)
	} else {
		log.Printf("Successfully encoded command into %d null-width characters.", len(nullWidthEncodedCommand))
	}

	// 4. Embed into cover text
	log.Println("Embedding encoded command into cover text...")
	textWithHiddenCommand := stealth.EmbedInText(coverText, nullWidthEncodedCommand)
	log.Println("Embedding complete.")

	// 5. Print the result for the user
	fmt.Println("\n--- Text with Embedded Command (for client use) ---")
	fmt.Println(textWithHiddenCommand)
	fmt.Println("--- End of Text --- (quoted version below for exact copying) --- ")
	fmt.Printf("%q\n", textWithHiddenCommand)

	log.Println("\nUser Instructions: Copy the text block above (or the quoted string). This text can now be placed on a webpage or social media post for the client to retrieve and execute the command.")

	// Perform server-side verification for immediate feedback
	log.Println("\n--- Server-Side Verification --- ")
	log.Println("Attempting to extract and decode the command from the generated text...")
	extractedBits := stealth.ExtractFromText(textWithHiddenCommand)
	if extractedBits == "" {
		log.Println("Verification failed: Could not extract any null-width bits.")
	} else {
		log.Printf("Verification: Extracted %d null-width bits.", len(extractedBits))
		decodedCommand, err := stealth.DecodeFromNullWidth(extractedBits)
		if err != nil {
			log.Printf("Verification failed: Error decoding extracted bits: %v", err)
			log.Printf("Verification: Extracted bits (raw) that failed decoding: %q", extractedBits)
		} else {
			if decodedCommand == command {
				log.Printf("Verification successful: Decoded command matches original: '%s'", decodedCommand)
			} else {
				log.Printf("Verification failed: Mismatch! Original: '%s', Decoded: '%s'", command, decodedCommand)
				log.Printf("Verification: Original encoded bits: %q", nullWidthEncodedCommand)
				log.Printf("Verification: Extracted bits (raw): %q", extractedBits)
			}
		}
	}
	log.Println("--- End of Server-Side Verification --- ")
	log.Println("Server finished operation.")
}
