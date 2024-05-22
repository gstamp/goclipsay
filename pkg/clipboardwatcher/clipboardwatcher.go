package clipboardwatcher

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gstamp/goclipsay/pkg/tts"
	"golang.design/x/clipboard"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type ClipboardWatcher struct {
	tts                                 tts.TTS
	mode                                string
	readEnglishTextWithJapaneseLanguage bool
}

func NewClipboardWatcher(tts tts.TTS, mode string, readEnglishTextWithJapaneseLanguage bool) *ClipboardWatcher {
	if mode != "en" && mode != "jp" {
		panic("language must be either en or jp")
	}

	return &ClipboardWatcher{
		tts:                                 tts,
		mode:                                mode,
		readEnglishTextWithJapaneseLanguage: readEnglishTextWithJapaneseLanguage,
	}
}

func (cw *ClipboardWatcher) Watch() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.Background(), clipboard.FmtText)

	for {
		for data := range ch {
			text := string(data)

			if cw.mode == "en" || (cw.mode == "jp" && isJapanese(text) && !cw.readEnglishTextWithJapaneseLanguage) {
				err = cw.handleClipboardChange(text)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}

func (cw *ClipboardWatcher) handleClipboardChange(data string) error {
	// replace carriage return and line feed with blank
	data = strings.ReplaceAll(data, "\r", "")
	data = strings.ReplaceAll(data, "\n", "")

	println(data)

	// Call Azure Speech API
	// https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/rest-text-to-speech
	body, err := cw.tts.RequestTTS(data, cw.mode == "en")
	if err != nil {
		return err
	}
	defer body.Close()

	streamer, format, err := mp3.Decode(body)
	if err != nil {
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	fmt.Println("Playing...")

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done

	return nil
}

func isJapanese(s string) bool {
	pattern := "[\u3040-\u309F\u30A0-\u30FF\u4E00-\u9FFF\uFF70-\uFFEF]+"
	match, _ := regexp.MatchString(pattern, s)

	return match
}
