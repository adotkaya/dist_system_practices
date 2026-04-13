package main

/*
  Task 1.1: Basic JSON Message Parser
  Reads JSON messages from stdin line-by-line, parses them, and extracts fields.
  Prints "PARSED: src|dest|body_type" to stdout for validation.
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Message represents a Maelstrom message
type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
			continue
		}

		// Extract body type with default fallback
		bodyType := "unknown"
		if typeVal, ok := msg.Body["type"]; ok {
			if typeStr, ok := typeVal.(string); ok {
				bodyType = typeStr
			}
		}

		// THIS IS THE KEY PART - Print to stdout for validation
		fmt.Printf("PARSED: %s|%s|%s\n", msg.Src, msg.Dest, bodyType)

		// Optional: Log to stderr for debugging
		fmt.Fprintf(os.Stderr, "Received from %s to %s: %v\n", msg.Src, msg.Dest, msg.Body)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Scanner error:", err)
	}
}
