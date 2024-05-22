package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OpenAITTSClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewOpenAIClient() (*OpenAITTSClient, error) {
	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("environment variable OPENAI_KEY not present")
	}

	return &OpenAITTSClient{
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
	}, nil
}

func (client *OpenAITTSClient) RequestTTS(clip string, englishMode bool) (io.ReadCloser, error) {

	payload := map[string]interface{}{
		"input": clip,
		"model": "tts-1-hd", // Use tts-1 or tts-1-hd based on your needs
		"voice": "onyx",     // Choose a voice preset
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/speech", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)

	// Make the HTTP request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (client *OpenAITTSClient) Type() string {
	return "OpenAI"
}
