package tts

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type AzureTTSClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewAzureClient() (*AzureTTSClient, error) {
	apiKey := os.Getenv("AZURE_SPEECH_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("environment variable AZURE_SPEECH_KEY not present")
	}

	return &AzureTTSClient{
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
	}, nil
}

func (client *AzureTTSClient) RequestTTS(clip string, englishMode bool) (io.ReadCloser, error) {
	token, err := client.requestToken()
	if err != nil {
		return nil, err
	}

	voice := ""
	if englishMode {
		voice = "Microsoft Server Speech Text to Speech Voice (en-AU, NatashaNeural)"
	} else {
		voice = "Microsoft Server Speech Text to Speech Voice (ja-JP, KeitaNeural)"
	}
	body := fmt.Sprintf("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='en-US'><voice name='%s'>%s</voice></speak>", voice, clip)
	bearer := strings.Join([]string{"Bearer ", *token}, "")
	url := "https://australiaeast.tts.speech.microsoft.com/cognitiveservices/v1"
	return client.post(url, map[string]string{
		"X-Microsoft-OutputFormat": "audio-16khz-64kbitrate-mono-mp3",
		"User-Agent":               "australiaeast",
		"Content-Type":             "application/ssml+xml",
		"Authorization":            bearer,
	}, body)
}

func (client *AzureTTSClient) Type() string {
	return "Azure"
}

func (client *AzureTTSClient) requestToken() (*string, error) {
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": client.apiKey,
		"Content-Length":            "0",
	}

	body, err := client.post(
		"https://australiaeast.api.cognitive.microsoft.com/sts/v1.0/issuetoken",
		headers,
		"",
	)
	if err != nil {
		return nil, fmt.Errorf("requesting TTS from azure: %w", err)
	}
	defer body.Close()

	tokenBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	token := string(tokenBytes)
	return &token, nil
}

func (c *AzureTTSClient) post(url string, headers map[string]string, body string) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	return resp.Body, err
}
