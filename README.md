# GoNull C2 Framework

A C2 (Command and Control) Remote Access Tool (RAT) system that uses null-width characters hidden on public webpages to send commands.

**Disclaimer:** This project is for educational and demonstration purposes only, specifically for a university cybersecurity club. It should not be used for any malicious activities.

## Project Goals

*   Demonstrate steganography techniques for covert communication.
*   Explore C2 mechanisms.
*   Provide a learning platform for cybersecurity concepts.

## Features (Planned)

*   **Server:**
    *   Encodes commands using null-width characters.
    *   (Potentially) Serves or helps place text with embedded commands.
    *   Receives exfiltrated data from clients.
*   **Client (RAT):**
    *   Fetches web content.
    *   Extracts and decodes commands from null-width characters.
    *   Executes commands.
    *   Sends results back to the server.

## How it Works

The system relies on embedding binary representations of commands into regular text using Unicode null-width characters. For example:

*   `U+200B` (Zero Width Space) might represent a binary '0'.
*   `U+200C` (Zero Width Non-Joiner) might represent a binary '1'.

The client fetches a seemingly innocuous webpage, extracts these hidden characters, reconstructs the command, executes it, and sends back the output.

## Setup & Usage

This project requires Go to be installed on your system.

### 1. Clone the Repository (if applicable)

```bash
git clone <repository-url>
cd GoNull
```

### 2. Running the Server (`cmd/server/main.go`)

The server application is responsible for encoding a command you provide into null-width characters and embedding it within a piece of cover text.

**Steps:**

1.  Navigate to the project's root directory in your terminal.
2.  Run the server:
    ```bash
    go run cmd/server/main.go
    ```
3.  The application will prompt you to:
    *   `Enter command to encode:` Type the command you want the client to execute (e.g., `ls -la`, `whoami`).
    *   `Enter cover text (e.g., a sentence or paragraph):` Provide some text that will be used to hide the encoded command. If you leave this blank, a default placeholder text will be used.

4.  The server will then output:
    *   The cover text with the embedded null-width command.
    *   A quoted version of this text for easy and accurate copying.
    *   A server-side verification of whether the command could be extracted and decoded correctly from its own output.

**Example Server Interaction:**

```
$ go run cmd/server/main.go
Enter command to encode: id
Enter cover text (e.g., a sentence or paragraph): This is a friendly message for the club.

--- Text with Embedded Command ---
This is ‎​‌‌​‌​‌​​‌‌​‌‌​‌‏a friendly message for the club.
--- End of Text --- reproducing it exactly for verification ---
"This is \u200e\u200b\u200c\u200c\u200b\u200c\u200b\u200c\u200b\u200b\u200c\u200c\u200b\u200c\u200c\u200b\u200c\u200fa friendly message for the club."
...
2023/10/27 10:00:00 Server-side verification: Successfully decoded command: 'id'
```

5.  **Copy the output text** (the line between `--- Text with Embedded Command ---` and `--- End of Text ---`, or the quoted string). This text is what you would place on a public webpage, social media post, etc., for the client to retrieve.

### 3. Running the Client (`cmd/client/main.go`)

The client application is responsible for fetching or receiving text, extracting the hidden null-width command, decoding it, and (eventually) executing it.

**Steps:**

1.  Navigate to the project's root directory in your terminal.
2.  Run the client:
    ```bash
    go run cmd/client/main.go
    ```
3.  The application will prompt you to:
    *   `Paste the text containing the hidden command ... (End input with an empty line or Ctrl+D)`

4.  **Paste the text** you copied from the server's output (or from the webpage/social media where it was placed).
    *   If the text is multi-line, paste it and then press Enter on a new empty line to signal the end of input.
    *   Alternatively, you can use Ctrl+D after pasting.

5.  The client will then:
    *   Attempt to extract and decode the command.
    *   Print the decoded command.

**Example Client Interaction:**

```
$ go run cmd/client/main.go
Paste the text containing the hidden command (e.g., copied from a webpage or server output):
(End input with an empty line or Ctrl+D)
This is ‎​‌‌​‌​‌​​‌‌​‌‌​‌‏a friendly message for the club.
(You press Enter on an empty line here if pasting multi-line, or just Ctrl+D)

--- Decoded Command ---
Command: 'id'
--- End of Command --- 
...
2023/10/27 10:05:00 Next step: Implement command execution for the decoded command.
```

### Next Steps for the Project

*   Implement actual command execution on the client-side.
*   Develop methods for the client to fetch text from live URLs instead of manual pasting.
*   Explore more sophisticated embedding and extraction techniques.
*   (Potentially) Implement a way for the client to send back command output to the server (though this is a significant extension and might also use steganography or other covert channels).

## Contributing

(To be added)