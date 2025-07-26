package clock

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Seven segment display segment assignments:
//
//   --a--
//  |     |
//  f     b
//  |     |
//   --g--
//  |     |
//  e     c
//  |     |
//   --d--

type SegmentRect struct {
	Position fyne.Position
	Size     fyne.Size
}

type SevenSegmentDisplay struct {
	Segments      [10][7]bool
	SegmentShapes [7]SegmentRect
	onColor       color.Color
	offColor      color.Color
	strokeColor   color.Color
	x, y          int
}

type DigitalClock struct {
	SevenSegmentDisplay *SevenSegmentDisplay
	ClockFace           *fyne.Container
	digits              []*fyne.Container
	digitWidth          int
	digitSpacing        int
	mode24hr            bool
}

func NewSevenSegmentDisplay(onColor, offColor, strokeColor color.Color) *SevenSegmentDisplay {
	segmentShapes := newDigitalSegmentShapes()
	segments := newDigitalSegmentBoolMap()

	return &SevenSegmentDisplay{
		Segments:      segments.Segments,
		SegmentShapes: segmentShapes,
		onColor:       onColor,
		offColor:      offColor,
		strokeColor:   strokeColor,
		x:             0,
		y:             0,
	}
}

func NewDigitalClock(mode24hr bool, onColor, offColor, strokeColor color.Color, digitalWidth, digitalSpacing int) *DigitalClock {
	time := NewTickData()

	ssd := NewSevenSegmentDisplay(onColor, offColor, strokeColor)

	HrTensDigit, HrOnesDigit := 0, 0

	if mode24hr {
		HrTensDigit = time.Hr24TensDigit
		HrOnesDigit = time.Hr24OnesDigit
	} else {
		HrTensDigit = time.Hr12TensDigit
		HrOnesDigit = time.Hr12OnesDigit
	}

	colon1 := drawColon(onColor)
	colon2 := drawColon(onColor)

	digitsContainer := []*fyne.Container{
		ssd.drawDigit(HrTensDigit),
		ssd.drawDigit(HrOnesDigit),
		container.NewWithoutLayout(colon1),
		ssd.drawDigit(time.MinTensDigit),
		ssd.drawDigit(time.MinOnesDigit),
		container.NewWithoutLayout(colon2),
		ssd.drawDigit(time.SecTensDigit),
		ssd.drawDigit(time.SecOnesDigit),
	}

	x := float32(0)
	for i, digitObj := range digitsContainer {
		var w float32
		if i == 2 || i == 5 {
			w = float32(digitalWidth) * 0.2 // or set a colonWidth variable
		} else {
			w = float32(digitalWidth)
		}

		digitObj.Resize(fyne.NewSize(w, digitObj.Size().Height))
		digitObj.Move(fyne.NewPos(x, 0))

		x += w + float32(digitalSpacing)
	}

	ClockFace := container.NewWithoutLayout()
	for _, digit := range digitsContainer {
		ClockFace.Add(digit)
	}

	return &DigitalClock{
		ClockFace:           ClockFace,
		digits:              digitsContainer,
		SevenSegmentDisplay: ssd,
		digitWidth:          digitalWidth,
		digitSpacing:        digitalSpacing,
		mode24hr:            mode24hr,
	}
}

func newDigitalSegmentShapes() [7]SegmentRect {
	return [7]SegmentRect{
		{ // a - top horizontal
			Position: fyne.NewPos(10, 0),
			Size:     fyne.NewSize(50, 10),
		},
		{ // b - top-right vertical
			Position: fyne.NewPos(60, 10),
			Size:     fyne.NewSize(10, 40),
		},
		{ // c - bottom-right vertical
			Position: fyne.NewPos(60, 60),
			Size:     fyne.NewSize(10, 40),
		},
		{ // d - bottom horizontal
			Position: fyne.NewPos(10, 100),
			Size:     fyne.NewSize(50, 10),
		},
		{ // e - bottom-left vertical
			Position: fyne.NewPos(0, 60),
			Size:     fyne.NewSize(10, 40),
		},
		{ // f - top-left vertical
			Position: fyne.NewPos(0, 10),
			Size:     fyne.NewSize(10, 40),
		},
		{ // g - middle horizontal
			Position: fyne.NewPos(10, 50),
			Size:     fyne.NewSize(50, 10),
		},
	}
}

func newDigitalSegmentBoolMap() *SevenSegmentDisplay {
	return &SevenSegmentDisplay{
		Segments: [10][7]bool{
			//  a, b, c, d, e, f, g
			{true, true, true, true, true, true, false},     // 0
			{false, true, true, false, false, false, false}, // 1
			{true, true, false, true, true, false, true},    // 2
			{true, true, true, true, false, false, true},    // 3
			{false, true, true, false, false, true, true},   // 4
			{true, false, true, true, false, true, true},    // 5
			{true, false, true, true, true, true, true},     // 6
			{true, true, true, false, false, false, false},  // 7
			{true, true, true, true, true, true, true},      // 8
			{true, true, true, true, false, true, true},     // 9
		},
		SegmentShapes: newDigitalSegmentShapes(),
	}
}

func (ssd *SevenSegmentDisplay) drawDigit(digit int) *fyne.Container {
	segments := []fyne.CanvasObject{}

	for i, seg := range ssd.SegmentShapes {
		segColor := ssd.offColor
		strokeColor := ssd.offColor

		if ssd.Segments[digit][i] {
			segColor = ssd.onColor
			strokeColor = ssd.strokeColor
		}

		rect := canvas.NewRectangle(segColor)
		rect.StrokeWidth = 2
		rect.StrokeColor = strokeColor
		rect.Resize(seg.Size)
		//apply offset
		rect.Move(fyne.NewPos(seg.Position.X+float32(ssd.x), seg.Position.Y+float32(ssd.y)))

		segments = append(segments, rect)
	}

	return container.NewWithoutLayout(segments...)
}

func drawColon(onColor color.Color) *fyne.Container {
	topDot := canvas.NewRectangle(onColor)
	topDot.Resize(fyne.NewSize(10, 10))

	bottomDot := canvas.NewRectangle(onColor)
	bottomDot.Resize(fyne.NewSize(10, 10))

	colon := container.NewWithoutLayout(topDot, bottomDot)
	topDot.Move(fyne.NewPos(0, 20))
	bottomDot.Move(fyne.NewPos(0, 60))

	return colon
}

func (d *DigitalClock) Update(t *TickData) {

	HrTensDigit, HrOnesDigit := 0, 0

	if d.mode24hr {
		HrTensDigit = t.Hr24TensDigit
		HrOnesDigit = t.Hr24OnesDigit
	} else {
		HrTensDigit = t.Hr12TensDigit
		HrOnesDigit = t.Hr12OnesDigit
	}

	newDigits := []*fyne.Container{
		d.SevenSegmentDisplay.drawDigit(HrTensDigit),
		d.SevenSegmentDisplay.drawDigit(HrOnesDigit),
		d.SevenSegmentDisplay.drawDigit(t.MinTensDigit),
		d.SevenSegmentDisplay.drawDigit(t.MinOnesDigit),
		d.SevenSegmentDisplay.drawDigit(t.SecTensDigit),
		d.SevenSegmentDisplay.drawDigit(t.SecOnesDigit),
	}

	// Only refresh digits
	digitIndex := 0
	for _, obj := range d.ClockFace.Objects {
		container, ok := obj.(*fyne.Container)
		if ok && len(container.Objects) == 7 {
			container.Objects = newDigits[digitIndex].Objects
			container.Refresh()
			digitIndex++
		}
	}
}
