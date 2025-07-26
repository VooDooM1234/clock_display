package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"temp.com/go-clock/clock"
)

func main() {
	a := app.New()
	w := a.NewWindow("Its Clocking time!")
	t := clock.NewTickData()
	w.Resize(fyne.NewSize(800, 600)) // wider for GIF

	const cx, cy, radius = 100, 100, 80
	const numberOfClocks = 3

	onColor := color.RGBA{R: 255, G: 150, B: 50, A: 255}
	offColor := color.RGBA{R: 50, G: 50, B: 50, A: 255}
	strokeColor := color.RGBA{R: 200, G: 100, B: 30, A: 255}

	ringClockSecColor := color.RGBA{R: 0, G: 255, B: 100, A: 255}
	ringClockMinColor := color.RGBA{R: 50, G: 150, B: 255, A: 255}
	ringClockHrColor := color.RGBA{R: 255, G: 80, B: 80, A: 255}

	const cxRing, cyRing, radiusRing = 100, 100, 80
	const digitalWidth, digitalSpacing = 70, 10

	analogClock := clock.NewAnalogClock(cx, cy, radius)
	digitalClock := clock.NewDigitalClock(true, onColor, offColor, strokeColor, digitalWidth, digitalSpacing)
	ringClock := clock.NewRingClock(cxRing, cyRing, radiusRing, ringClockSecColor, ringClockMinColor, ringClockHrColor, offColor, strokeColor)
	ringClock.BackFillArcsContainer()

	analogClockContainer := container.NewWithoutLayout(
		analogClock.ClockFace,
		analogClock.HourHand,
		analogClock.MinuteHand,
		analogClock.SecondHand,
	)

	digitalClockContainer := container.NewWithoutLayout(
		digitalClock.ClockFace,
	)

	ringClockContainer := container.NewWithoutLayout(
		ringClock.ClockFace,
	)

	clocks := container.NewGridWithRows(numberOfClocks,
		analogClockContainer,
		digitalClockContainer,
		ringClockContainer,
	)

	datBoi := NewAnimatedGIF("local/images/Dat_boi.gif")
	datBoi.Start() // start animation

	content := container.NewGridWithColumns(2,
		clocks,
		datBoi.Image, // just add the *canvas.Image
	)

	w.SetContent(content)

	// clock updater
	go func() {
		for range time.Tick(time.Second) {
			fyne.Do(func() {
				t.Update()
				analogClock.Update(t)
				digitalClock.Update(t)
				ringClock.Update(t)
			})
		}
	}()

	w.ShowAndRun()
}
