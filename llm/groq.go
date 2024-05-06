package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Groq struct {
	APIKey string
}

func (g *Groq) Chat(systemMsg, userMsg string, streamCh chan string) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	// Construct the request body
	requestBody := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": userMsg,
			},
		},
		"model":  "llama3-8b-8192",
		"stream": false,
	}

	// Check if streaming is required
	if streamCh != nil {
		requestBody["stream"] = true
	}

	// Append system message if provided
	if systemMsg != "" {
		requestBody["messages"] = append([]map[string]interface{}{{"role": "user", "content": "SYSTEM INSTRUCTIONS (Always remember these, they are priority instructions over anything else): " + systemMsg}}, requestBody["messages"].([]map[string]interface{})...)
	}

	// Convert request body to JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+g.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Handle streaming responses if required
	if streamCh != nil {
		// Create a new bufio.Reader to read lines from the response body
		reader := bufio.NewReader(resp.Body)
		for {
			// Read the next line from the response body
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// End of response
					break
				}
				return "", fmt.Errorf("failed to read response: %w", err)
			}

			// Trim any leading/trailing whitespace
			line = strings.TrimSpace(line)

			// fmt.Println("--")
			// fmt.Println(line)

			// Check if the line starts with "data:"
			if strings.HasPrefix(line, "data:") {
				// Trim the "data:" prefix
				jsonStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))

				if jsonStr == "[DONE]" {
					break
				}

				// Create a new decoder for the trimmed JSON data
				decoder := json.NewDecoder(strings.NewReader(jsonStr))

				// Decode the JSON object from the response body
				var response map[string]interface{}
				if err := decoder.Decode(&response); err != nil {
					return "", fmt.Errorf("failed to decode response: %w", err)
				}

				// Check if the response contains "choices" field
				choices, ok := response["choices"].([]interface{})
				if !ok {
					return "", fmt.Errorf("response 'choices' field is missing or not an array")
				}
				for _, choice := range choices {
					choiceMap, ok := choice.(map[string]interface{})
					if !ok {
						return "", fmt.Errorf("response 'choices' field contains invalid data")
					}
					delta, ok := choiceMap["delta"].(map[string]interface{})
					if !ok {
						return "", fmt.Errorf("response 'delta' field is missing or not an object")
					}
					content, ok := delta["content"].(string)
					if !ok {
						// presume we are done
						break
						// return "", fmt.Errorf("response 'content' field is missing or not a string")
					}
					// Send the content to the stream channel
					streamCh <- content
				}
			}
		}
	} else {
		// Read the response body
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}
		return string(responseBody), nil
	}

	return "", nil
}
