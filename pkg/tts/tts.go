package tts

import "io"

type TTS interface {
	RequestTTS(clip string, englishMode bool) (io.ReadCloser, error)
	Type() string
}
