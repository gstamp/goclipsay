package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gstamp/goclipsay/pkg/clipboardwatcher"
	"github.com/gstamp/goclipsay/pkg/tts"
)

func autoDetectTTS() (tts.TTS, error) {
	if os.Getenv("OPENAI_KEY") != "" {
		return tts.NewOpenAIClient()
	} else if os.Getenv("AZURE_SPEECH_KEY") != "" {
		return tts.NewAzureClient()
	} else {
		return nil, fmt.Errorf("no tts service found, set either OPENAI_KEY or AZURE_SPEECH_KEY environment variable")
	}
}

func main() {
	azure := flag.Bool("m", false, "use microsoft azure tts service")
	openai := flag.Bool("a", false, "use openai tts service")
	mode := flag.String("l", "jp", "language mode, en or jp - does not matter for openai tts service")
	readEnglishTextWithJapaneseLanguage := flag.Bool("e", false, "if lang is jp, read english text with japanese language voice")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}
	if *mode != "en" && *mode != "jp" {
		panic("language mode must be either en or jp")
	}

	var ttsClient tts.TTS
	var err error

	if *azure {
		ttsClient, err = tts.NewAzureClient()
		if err != nil {
			panic(err)
		}
	} else if *openai {
		ttsClient, err = tts.NewOpenAIClient()
		if err != nil {
			panic(err)
		}
	} else if ttsClient, err = autoDetectTTS(); err != nil {
		flag.PrintDefaults()
		return
	}

	fmt.Println("Using TTS service: ", ttsClient.Type())

	watcher := clipboardwatcher.NewClipboardWatcher(ttsClient, *mode, *readEnglishTextWithJapaneseLanguage)

	watcher.Watch()
}
