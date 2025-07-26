package main

import (
	"image/gif"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// AnimatedGIF holds the GIF frames + the image object
type AnimatedGIF struct {
	Image  *canvas.Image
	frames []*canvas.Image
	delays []time.Duration
}

// NewAnimatedGIF loads a GIF from file and prepares it
func NewAnimatedGIF(path string) *AnimatedGIF {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	g, err := gif.DecodeAll(file)
	if err != nil {
		panic(err)
	}

	var frames []*canvas.Image
	var delays []time.Duration

	for i, img := range g.Image {
		frames = append(frames, canvas.NewImageFromImage(img))
		delays = append(delays, time.Duration(g.Delay[i])*10*time.Millisecond)
	}

	initial := canvas.NewImageFromImage(g.Image[0])
	initial.FillMode = canvas.ImageFillOriginal

	return &AnimatedGIF{
		Image:  initial,
		frames: frames,
		delays: delays,
	}
}

// Start begins animating the GIF in a loop
func (a *AnimatedGIF) Start() {
	go func() {
		frame := 0
		for {
			time.Sleep(a.delays[frame])
			frame = (frame + 1) % len(a.frames)

			fyne.Do(func() {
				// swap frame image
				a.Image.Image = a.frames[frame].Image
				a.Image.Refresh()
			})
		}
	}()
}
