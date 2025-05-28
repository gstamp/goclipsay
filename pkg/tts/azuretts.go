package tts

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	maxRetries    = 3
	retryInterval = 2 * time.Second
	burstSize     = 1
	speechUrl     = "https://eastus.tts.speech.microsoft.com/cognitiveservices/v1"
	issueTokenUrl = "https://eastus.api.cognitive.microsoft.com/sts/v1.0/issuetoken"
	englishVoice  = "Microsoft Server Speech Text to Speech Voice (en-AU, NatashaNeural)"
	japaneseVoice = "ja-JP-Masaru:DragonHDLatestNeural"
	userAgent     = "goclipsay"

	// Australian server endpoints
	// speechUrl     = "https://australiaeast.tts.speech.microsoft.com/cognitiveservices/v1"
	// issueTokenUrl = "https://australiaeast.api.cognitive.microsoft.com/sts/v1.0/issuetoken"

)

// Azure free tier allows 20 requests per minute
var rateLimit = rate.Every(3 * time.Second) // One request every 3 seconds

type AzureTTSClient struct {
	httpClient *http.Client
	apiKey     string
	limiter    *rate.Limiter
}

func NewAzureClient() (*AzureTTSClient, error) {
	apiKey := os.Getenv("AZURE_SPEECH_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("environment variable AZURE_SPEECH_KEY not present")
	}

	return &AzureTTSClient{
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
		limiter:    rate.NewLimiter(rateLimit, burstSize),
	}, nil
}

func (client *AzureTTSClient) RequestTTS(clip string, englishMode bool) (io.ReadCloser, error) {
	token, err := client.requestToken()
	if err != nil {
		return nil, err
	}

	voice := ""
	if englishMode {
		voice = englishVoice
	} else {
		voice = japaneseVoice
	}
	body := fmt.Sprintf("<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='ja-JP'><voice name='%s'>%s</voice></speak>",
		html.EscapeString(voice), html.EscapeString(clip))
	bearer := strings.Join([]string{"Bearer ", *token}, "")
	url := speechUrl
	return client.post(url, map[string]string{
		"X-Microsoft-OutputFormat": "audio-16khz-64kbitrate-mono-mp3",
		"User-Agent":               userAgent,
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
		issueTokenUrl,
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
	// Wait for rate limiter
	err := c.limiter.Wait(context.Background())
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	var resp *http.Response
	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", url, strings.NewReader(body))
		if err != nil {
			return nil, err
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			break
		}

		if attempt < maxRetries {
			resp.Body.Close()
			time.Sleep(retryInterval)
			continue
		}
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected response status while generating TTS: %d %s", resp.StatusCode, resp.Status)
	}

	return resp.Body, nil
}
