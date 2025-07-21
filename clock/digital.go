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
	ClockFace           fyne.CanvasObject
	digits              []fyne.CanvasObject
	digitWidth          int
	digitSpacing        int
	HrOnesDigit         int
	HrTensDigit         int
	MinOnesDigit        int
	MinTensDigit        int
	SecOnesDigit        int
	SecTensDigit        int
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

func NewDigitalClock(onColor, offColor, strokeColor color.Color, digitalWidth, digitalSpacing int) *DigitalClock {
	time := NewTickData()

	ssd := NewSevenSegmentDisplay(onColor, offColor, strokeColor)
	ssd.SegmentShapes = newDigitalSegmentShapes()
	ssd.Segments = newDigitalSegmentBoolMap().Segments
	ssd.onColor = onColor
	ssd.offColor = offColor
	ssd.strokeColor = strokeColor

	HrOnesDigit := ssd.drawDigit(time.HrOnesDigit)
	HrTensDigit := ssd.drawDigit(time.HrTensDigit)
	MinOnesDigit := ssd.drawDigit(time.MinOnesDigit)
	MinTensDigit := ssd.drawDigit(time.MinTensDigit)
	SecOnesDigit := ssd.drawDigit(time.SecOnesDigit)
	SecTensDigit := ssd.drawDigit(time.SecTensDigit)

	digits := []fyne.CanvasObject{
		HrTensDigit,
		HrOnesDigit,
		MinTensDigit,
		MinOnesDigit,
		SecTensDigit,
		SecOnesDigit,
	}

	for i, digitObj := range digits {
		digitObj.Move(fyne.NewPos(float32(i*(digitalWidth+digitalSpacing)), 0))
	}

	ClockFace := container.NewWithoutLayout(digits...)

	return &DigitalClock{
		ClockFace:           ClockFace,
		digits:              digits,
		SevenSegmentDisplay: NewSevenSegmentDisplay(onColor, offColor, strokeColor),
		digitWidth:          digitalWidth,
		digitSpacing:        digitalSpacing,
		HrOnesDigit:         0,
		HrTensDigit:         0,
		MinOnesDigit:        0,
		MinTensDigit:        0,
		SecOnesDigit:        0,
		SecTensDigit:        0,
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

func (ssd *SevenSegmentDisplay) drawDigit(digit int) fyne.CanvasObject {
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

func (d *DigitalClock) Update(t *TickData) {

	digitObjects := []fyne.CanvasObject{
		d.SevenSegmentDisplay.drawDigit(t.HrTensDigit),
		d.SevenSegmentDisplay.drawDigit(t.HrOnesDigit),
		d.SevenSegmentDisplay.drawDigit(t.MinTensDigit),
		d.SevenSegmentDisplay.drawDigit(t.MinOnesDigit),
		d.SevenSegmentDisplay.drawDigit(t.SecTensDigit),
		d.SevenSegmentDisplay.drawDigit(t.SecOnesDigit),
	}

	for i, newDigit := range digitObjects {
		existingContainer := d.digits[i].(*fyne.Container)

		// Clear old children
		existingContainer.Objects = nil

		// Add new children one by one
		newDigitContainer := newDigit.(*fyne.Container)
		for _, obj := range newDigitContainer.Objects {
			existingContainer.Add(obj)
		}
		existingContainer.Refresh()
	}
}
