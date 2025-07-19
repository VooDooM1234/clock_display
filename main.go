package main

import (
	"image/color"
	"math"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func drawClockFace(cx int, cy int, radius int) fyne.CanvasObject {

	hourMakerColor := color.Black
	hour_marker_multiplier := 0.8

	face := container.NewWithoutLayout()

	circle := canvas.NewCircle(color.White)
	circle.StrokeColor = color.Gray{Y: 0x99}
	circle.StrokeWidth = 5
	circle.Resize(fyne.NewSize(float32(radius*2), float32(radius*2)))
	circle.Move(fyne.NewPos(float32(cx-radius), float32(cy-radius)))

	face.Add(circle)

	// Hour markers

	hourMarkers := container.NewWithoutLayout()

	for i := 1; i <= 12; i++ {
		hourMarker := canvas.NewText(strconv.Itoa(i), hourMakerColor)
		hourMarker.TextSize = 18
		hourMarker.Refresh()
		textSize := hourMarker.MinSize()
		angle := float64(i%12) * (math.Pi / 6)

		x := float64(cx) + float64(radius)*hour_marker_multiplier*float64(math.Sin(angle)) - float64(textSize.Width)/2
		y := float64(cy) - float64(radius)*hour_marker_multiplier*float64(math.Cos(angle)) - (float64(textSize.Height) / 2) - 6

		hourMarker.Move(fyne.NewPos(float32(x), float32(y)))
		hourMarkers.Add(hourMarker)
	}

	return container.NewVBox(face, hourMarkers)
}

func drawClockHand(cx int, cy int, x int, y int, handColor color.Color) *canvas.Line {
	hand := canvas.NewLine(handColor)
	hand.StrokeWidth = 2
	hand.Position1 = fyne.NewPos(float32(cx), float32(cy))
	hand.Position2 = fyne.NewPos(float32(x), float32(y))
	hand.Refresh()
	return hand
}

// @return x, y coordinates of the clock hand position
func getClockHandPosition(cx int, cy int, radius int, angle float64) (int, int) {

	angle_rad := angle * (math.Pi / 180)

	x := int(float64(cx) + float64(radius)*math.Sin(angle_rad))
	y := int(float64(cy) - float64(radius)*math.Cos(angle_rad))

	return x, y
}

func updateHand(hand *canvas.Line, cx int, cy int, radius int, angle float64) {

	X, Y := getClockHandPosition(cx, cy, radius-20, angle)
	hand.Position2 = fyne.NewPos(float32(X), float32(Y))
	canvas.Refresh(hand)
}

func main() {
	a := app.New()
	w := a.NewWindow("Its Clocking time!")

	// Set window size
	w.Resize(fyne.NewSize(400, 400))

	// Define circle parameters
	cx, cy, radius := 100, 100, 80

	clockFace := drawClockFace(cx, cy, radius)

	currentTime := time.Now()

	hour := (currentTime.Hour() + 1) % 12
	minute := currentTime.Minute()
	second := currentTime.Second()

	hourAngle := float64(hour) * 30
	minuteAngle := float64(minute) * 6
	secondAngle := float64(second) * 6

	hourX, hourY := getClockHandPosition(cx, cy, radius, hourAngle)
	minuteX, minuteY := getClockHandPosition(cx, cy, radius, minuteAngle)
	secondX, secondY := getClockHandPosition(cx, cy, radius, secondAngle)

	hourHand := drawClockHand(cx, cy, hourX, hourY, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	minuteHand := drawClockHand(cx, cy, minuteX, minuteY, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	secondHand := drawClockHand(cx, cy, secondX, secondY, color.RGBA{R: 0, G: 0, B: 255, A: 255})

	content := container.NewWithoutLayout()
	content.Add(clockFace)
	content.Add(hourHand)
	content.Add(minuteHand)
	content.Add(secondHand)

	w.SetContent(content)

	go func() {
		for range time.Tick(time.Second) {

			currentTime := time.Now()

			hour := (currentTime.Hour() + 1) % 12
			minute := currentTime.Minute()
			second := currentTime.Second()

			hourAngle := float64(hour) * 30
			minuteAngle := float64(minute) * 6
			secondAngle := float64(second) * 6

			fyne.Do(func() {
				updateHand(hourHand, cx, cy, radius, hourAngle)
				updateHand(minuteHand, cx, cy, radius, minuteAngle)
				updateHand(secondHand, cx, cy, radius, secondAngle)
			})
		}
	}()

	w.ShowAndRun()
}
