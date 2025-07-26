package clock

import (
	"image/color"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type AnalogClock struct {
	HourHand    *canvas.Line
	MinuteHand  *canvas.Line
	SecondHand  *canvas.Line
	ClockFace   fyne.CanvasObject
	cx, cy      int
	radius      int
	HourAngle   float64
	MinuteAngle float64
	SecondAngle float64
}

func NewAnalogClock(cx, cy, radius int) *AnalogClock {

	time := NewTickData()

	offset := 20

	hourAngle, minuteAngle, secondAngle := getClockHandAngles(time)

	clockFace := drawClockFace(cx, cy, radius)
	hourX, hourY := getClockHandPosition(cx, cy, radius-offset, hourAngle)
	minuteX, minuteY := getClockHandPosition(cx, cy, radius-offset, minuteAngle)
	secondX, secondY := getClockHandPosition(cx, cy, radius-offset, secondAngle)

	hourHand := drawClockHand(cx, cy, hourX, hourY, color.RGBA{255, 0, 0, 255})
	minuteHand := drawClockHand(cx, cy, minuteX, minuteY, color.RGBA{0, 255, 0, 255})
	secondHand := drawClockHand(cx, cy, secondX, secondY, color.RGBA{0, 0, 255, 255})

	return &AnalogClock{
		HourHand:    hourHand,
		MinuteHand:  minuteHand,
		SecondHand:  secondHand,
		ClockFace:   clockFace,
		cx:          cx,
		cy:          cy,
		radius:      radius,
		HourAngle:   hourAngle,
		MinuteAngle: minuteAngle,
		SecondAngle: secondAngle,
	}
}

func getClockHandAngles(time *TickData) (float64, float64, float64) {
	hourAngle := (float64(time.Hour12) + float64(time.Minute)/60.0) * 30.0
	minuteAngle := float64(time.Minute) * 6
	secondAngle := float64(time.Second) * 6

	return hourAngle, minuteAngle, secondAngle
}

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
	offset := 20
	X, Y := getClockHandPosition(cx, cy, radius-offset, angle)
	hand.Position2 = fyne.NewPos(float32(X), float32(Y))
	canvas.Refresh(hand)
}

func (a *AnalogClock) Update(t *TickData) {
	a.HourAngle, a.MinuteAngle, a.SecondAngle = getClockHandAngles(t)
	updateHand(a.HourHand, a.cx, a.cy, a.radius, a.HourAngle)
	updateHand(a.MinuteHand, a.cx, a.cy, a.radius, a.MinuteAngle)
	updateHand(a.SecondHand, a.cx, a.cy, a.radius, a.SecondAngle)
}
