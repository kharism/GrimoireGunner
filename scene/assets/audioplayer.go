package assets

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var musicPlayerCh chan *AudioPlayer
var audioContext *audio.Context

const (
	sampleRate = 48000
)

type AudioPlayer struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte

	volume128 int

	musicType musicType
	playSfx   bool
}
type musicType int

const (
	TypeOgg musicType = iota
	TypeMP3
)

// var AudioContext *audio.Context

func init() {
	audioContext = audio.CurrentContext()
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}

}

func (p *AudioPlayer) AudioPlayer() *audio.Player {
	return p.audioPlayer
}
func (p *AudioPlayer) Close() error {
	return p.audioPlayer.Close()
}
func (p *AudioPlayer) Update() error {
	select {
	case p.seBytes = <-p.seCh:
		// fmt.Println("SFX detected")
		// close(p.seCh)
		// p.playSfx = true

		// p.seCh = nil
	default:
	}
	p.PlaySEIfNeeded()
	return nil
}
func (p *AudioPlayer) ShouldPlaySE() bool {
	if p.seBytes == nil {
		// Bytes for the SE is not loaded yet.
		return false
	}
	// fmt.Println(p.seCh)
	return p.seCh != nil
}

func (p *AudioPlayer) PlaySEIfNeeded() {
	if !p.ShouldPlaySE() {
		return
	}
	sePlayer := p.audioContext.NewPlayerFromBytes(p.seBytes)
	sePlayer.Play()
	p.seBytes = nil
}
func (p *AudioPlayer) QueueSFXNoSampling(param []byte) {
	p.seCh <- param
}
func (p *AudioPlayer) QueueSFX(param []byte) {
	go func() {
		s, err := mp3.DecodeWithoutResampling(bytes.NewReader(param))
		// p.seCh <- param
		if err != nil {
			log.Fatal(err)
			return
		}
		b, err := io.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		p.seCh <- b
	}()

}
func NewAudioPlayer(audioByte []byte, musicType musicType) (*AudioPlayer, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	const bytesPerSample = 4 // TODO: This should be defined in audio package

	var s audioStream
	// audio, err := os.Open(audioPath)
	// if err != nil {
	// 	return nil, err
	// }
	// defer audio.Close()
	switch musicType {

	case TypeMP3:
		var err error
		s, err = mp3.DecodeWithoutResampling(bytes.NewReader(audioByte))
		if err != nil {
			return nil, err
		}
	default:
		panic("not reached")
	}

	p, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}

	player := &AudioPlayer{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128:    2,
		seCh:         make(chan []byte, 100),
		seBytes:      []byte{},
		musicType:    musicType,
	}
	if player.total == 0 {
		player.total = 1
	}

	// player.audioPlayer.Play()

	return player, nil
}
