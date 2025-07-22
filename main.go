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
	// t := clock.NewTickData()

	// Set window size
	w.Resize(fyne.NewSize(600, 600))

	// Define circle parameters
	const cx, cy, radius = 100, 100, 80

	onColor := color.RGBA{R: 255, G: 150, B: 50, A: 255}
	offColor := color.RGBA{R: 50, G: 50, B: 50, A: 255}
	strokeColor := color.RGBA{R: 200, G: 100, B: 30, A: 255}

	const digitalWidth, digitalSpacing = 70, 10

	analogClock := clock.NewAnalogClock(cx, cy, radius)
	digitalClock := clock.NewDigitalClock(true, onColor, offColor, strokeColor, digitalWidth, digitalSpacing)

	analogClockContainer := container.NewWithoutLayout(
		analogClock.ClockFace,
		analogClock.HourHand,
		analogClock.MinuteHand,
		analogClock.SecondHand,
	)

	digitalClockContainer := container.NewWithoutLayout(
		digitalClock.ClockFace,
	)

	content := container.NewGridWithRows(2,
		analogClockContainer,
		digitalClockContainer,
	)

	w.SetContent(content)

	go func() {
		for range time.Tick(time.Second) {
			t := clock.NewTickData()
			fyne.Do(func() {
				analogClock.Update(t)
				digitalClock.Update(t)
			})
		}
	}()

	w.ShowAndRun()
}
