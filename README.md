# GoClipSay

GoClipSay is a Go application that uses Text-to-Speech (TTS) services to read out text from the clipboard. It supports Microsoft Azure and OpenAI TTS services. 
It's main purpose is for reading Japanese text however it's possible to use if for English text as well.

## Usage

To run the application, use the following command:

```sh
go run main.go
```

You can specify the TTS service and language mode using command-line flags:

`-m`: Use Microsoft Azure TTS service.
`-a`: Use OpenAI TTS service.
`-l`: Language mode, either "en" for English or "jp" for Japanese. This does not matter for OpenAI TTS service. Japanese is the default.
`-j`: Only read Japanese text with Kana or Kanji characters. This is true by default.

For example, to use Azure TTS service and set the language mode to English, run:

```sh
go run main.go -m -l en
```

## Environment Variables

The application auto-detects the TTS service based on the following environment variables:

`OPENAI_KEY`: If set, the application uses OpenAI TTS service.
`AZURE_SPEECH_KEY`: If set, the application uses Azure TTS service.

If neither environment variable is set, the application will return an error.

Please note that use of these API's may incur a charge.

## Structure

The project is structured as follows:

- `main.go`: The entry point of the application.
- `pkg/clipboardwatcher/clipboardwatcher.go`: Contains the `ClipboardWatcher` struct and its methods.
- `pkg/tts/azuretts.go`: Contains the `AzureTTSClient` struct and its methods for interacting with Azure TTS service.
- `pkg/tts/openaitts.go`: Contains the `OpenAIClient` struct and its methods for interacting with OpenAI TTS service.
- `pkg/tts/tts.go`: Contains the `TTS` interface that both `AzureTTSClient` and `OpenAIClient` implement.

just some test changes to see what happens if you delete a branch and a PR exists on that branch
