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
    *   Helps place text with embedded commands (currently by printing to console).
    *   Receives exfiltrated data from clients (Future Goal).
*   **Client (RAT):**
    *   Fetches web content from a URL.
    *   Extracts and decodes commands from null-width characters found in the fetched content.
    *   Executes the decoded commands on the host system.
    *   Displays command output to the console.
    *   Sends results back to the server (Future Goal).

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

The client application is responsible for fetching web content from a URL, attempting to extract a hidden null-width command from that content, decoding it, and then executing it on the local system.

**Steps:**

1.  Navigate to the project's root directory in your terminal.
2.  Run the client:
    ```bash
    go run cmd/client/main.go
    ```
3.  The application will prompt you to:
    *   `Enter the URL to fetch content from:` Enter the full URL of the webpage where the command-carrying text is located (e.g., a Steam profile, a forum post).

4.  The client will then:
    *   Attempt to fetch the content from the provided URL.
    *   If successful, it will process the fetched content (e.g., HTML) to find and extract the null-width characters.
    *   Decode these characters back into a command string.
    *   Print the decoded command.
    *   Execute the decoded command and print its standard output/error.

**Example Client Interaction:**

Assume the server was used to embed the command `echo Hello from GoNull` into text on a webpage at `<URL_of_webpage_with_hidden_command>`.

```
$ go run cmd/client/main.go
Enter the URL to fetch content from: <URL_of_webpage_with_hidden_command>
YYYY/MM/DD HH:MM:SS Fetching content from: <URL_of_webpage_with_hidden_command>
YYYY/MM/DD HH:MM:SS Attempting to extract hidden command from fetched content...

--- Decoded Command ---
Command: 'echo Hello from GoNull'
--- End of Command --- 

YYYY/MM/DD HH:MM:SS Executing command: 'echo Hello from GoNull'

--- Command Output ---
Hello from GoNull
--- End of Output --- 

YYYY/MM/DD HH:MM:SS Client finished operation.
```

*(Note: The success of fetching and extraction depends on the accessibility of the URL, its content structure, and whether the null-width characters are preserved by the web platform. The command execution will run with the permissions of the user running the client application.)*

### Next Steps for the Project

*   Explore more sophisticated embedding and extraction techniques (e.g., dealing with specific HTML structures, character encodings, or dynamic content loaded by JavaScript).
*   Error handling and resilience for the client (e.g., retries, timeouts for HTTP requests, more robust command parsing).
*   Stealthier command execution (e.g., running in the background, avoiding console windows on some OSes - though this goes against the "no evasion" principle for this project, good to be aware of for general C2).
*   Implement a way for the client to send back command output to the server (e.g., via another steganographic message or a direct HTTP POST).

## Contributing

(To be added)