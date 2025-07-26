package clock

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"temp.com/go-clock/utils"
)

type RingClock struct {
	ClockFace                *fyne.Container   // Collection of all clock containers
	clocksContainer          []*fyne.Container //Collection ring of each ring, label and arcs grouping i.e. HH, MM, SS
	arcContainers            []*fyne.Container //Contains collection of arcs for filling the ring
	clockRings               []ClockRing
	cx, cy                   int
	radius                   int
	thickness                int
	spacing                  int
	strokeColor              color.Color
	numRings                 int
	offsetXarr               []int
	arcStrokeWidth           []float32
	timeAngles               []float64
	ringColorDimFactor       float32
	labelScaleFactorOfRadius float32
	arcOverShoot             int
	arcResolution            int
	arcNumResolutionLines    []int
}

// Specifc Data to each ring in th=e clock face
type ClockRing struct {
	Name       string
	epochLabel string
	X, Y       int
	onColor    color.Color
	offColor   color.Color
	labelColor color.Color
}

func NewRingClock(cx, cy, radius int, secondColor, MinuteColor, HourColor, offColor, strokeColor color.Color) *RingClock {
	time := NewTickData()

	const thickness, spacing = 40, 10
	const labelScaleFactorOfRadius = 0.3
	const ringColorDimFactor = 0.3
	const numRings = 3
	const arcOverShoot = 1
	const arcResolution = 1

	var arcNumResolutionLines []int
	var clocksContainers []*fyne.Container
	var arcContainers []*fyne.Container

	hourAngle, minuteAngle, secondAngle := getClockHandAngles(time)
	timeAngles := time.GetClockAngles()

	secondStrokeWidth := float32(radius) * (6 * math.Pi / 180)
	minuteStrokeWidth := secondStrokeWidth
	hourStrokeWidth := float32(radius) * (math.Pi / 6)

	arcStrokeWidth := []float32{secondStrokeWidth, minuteStrokeWidth, hourStrokeWidth}

	for _, width := range arcStrokeWidth {
		numLines := int(float32(math.Round(float64(width) / float64(arcResolution))))
		arcNumResolutionLines = append(arcNumResolutionLines, numLines)
	}

	offsetXarr := []int{}
	for i := 0; i < numRings; i++ {
		offsetX := i * ((radius * 2) + spacing)
		offsetXarr = append(offsetXarr, offsetX)
	}

	secondX, secondY := getClockHandPosition(offsetXarr[0], cy, radius, secondAngle)
	minuteX, minuteY := getClockHandPosition(offsetXarr[1], cy, radius, minuteAngle)
	hourX, hourY := getClockHandPosition(offsetXarr[2], cy, radius, hourAngle)

	clockRings := []ClockRing{
		{"Seconds", "SS", secondX, secondY, secondColor, utils.DimColor(secondColor, ringColorDimFactor), secondColor},
		{"Minutes", "MM", minuteX, minuteY, MinuteColor, utils.DimColor(MinuteColor, ringColorDimFactor), MinuteColor},
		{"Hours", "HH", hourX, hourY, HourColor, utils.DimColor(HourColor, ringColorDimFactor), HourColor},
	}

	//Construnct each temporal clock
	for i, r := range clockRings {
		ring := drawRing(cx+offsetXarr[i], cy, radius, thickness, theme.Color(theme.ColorNameBackground), clockRings[i].offColor)
		ringsContainers := container.NewWithoutLayout()
		ringsContainers.Add(ring)

		label := canvas.NewText(r.Name, clockRings[i].labelColor)
		label.TextSize = float32(radius) * labelScaleFactorOfRadius
		label.Alignment = fyne.TextAlignCenter
		label.TextStyle = fyne.TextStyle{Bold: true}
		label.Move(fyne.NewPos(float32(cx+offsetXarr[i]), float32(cy-(int(label.MinSize().Height)/2))))

		arcContainer := container.NewWithoutLayout()
		arcContainers = append(arcContainers, arcContainer)

		clockContainer := container.NewWithoutLayout(ring, label, arcContainer)

		clocksContainers = append(clocksContainers, clockContainer)
	}

	canvasObjects := make([]fyne.CanvasObject, len(clocksContainers))
	for i, c := range clocksContainers {
		canvasObjects[i] = c
	}

	ClockFace := container.NewWithoutLayout(canvasObjects...)

	return &RingClock{
		ClockFace:                ClockFace,
		clocksContainer:          clocksContainers,
		arcContainers:            arcContainers,
		cx:                       cx,
		cy:                       cy,
		radius:                   radius,
		thickness:                thickness,
		spacing:                  spacing,
		strokeColor:              strokeColor,
		timeAngles:               timeAngles,
		numRings:                 numRings,
		offsetXarr:               offsetXarr,
		arcStrokeWidth:           arcStrokeWidth,
		clockRings:               clockRings,
		ringColorDimFactor:       ringColorDimFactor,
		labelScaleFactorOfRadius: labelScaleFactorOfRadius,
		arcOverShoot:             arcOverShoot,
		arcResolution:            arcResolution,
		arcNumResolutionLines:    arcNumResolutionLines,
	}
}

func drawRing(cx, cy, radius, thickness int, innerColor, outerColor color.Color) *fyne.Container {

	outerRing := canvas.NewCircle(outerColor)
	innerRing := canvas.NewCircle(innerColor)

	outerRing.Resize(fyne.NewSize(float32(radius*2), float32(radius*2)))
	outerRing.Move(fyne.NewPos(float32(cx-radius), float32(cy-radius)))

	innerRing.Resize(fyne.NewSize(float32(radius*2-thickness), float32(radius*2-thickness)))
	innerRing.Move(fyne.NewPos(float32(cx-radius+(thickness/2)), float32(cy-radius+(thickness/2))))

	ring := container.NewWithoutLayout(
		outerRing,
		innerRing,
	)
	return ring
}

// Draw singular arc using fyne.canvas.Line
func (r *RingClock) drawArc(cx, cy, x, y int, stroke float32, arcColor color.Color) *canvas.Line {

	arc := canvas.NewLine(arcColor)
	arc.StrokeWidth = float32(r.arcResolution)
	arc.Position1 = fyne.NewPos(float32(cx), float32(cy))
	arc.Position2 = fyne.NewPos(float32(x), float32(y))
	arc.StrokeWidth = stroke

	return arc
}

func (r *RingClock) drawArcsToResolution(cx, cy, x, y int, stroke float32, arcColor color.Color) []*canvas.Line {
	var linesArr []*canvas.Line
	for i := range r.clockRings {
		numLines := r.arcNumResolutionLines[i]

		for j := 0; j < numLines; j++ {
			arc := canvas.NewLine(arcColor)
			arc.StrokeWidth = stroke
			arc.Position1 = fyne.NewPos(float32(cx), float32(cy))
			arc.Position2 = fyne.NewPos(float32(x), float32(y))
			linesArr = append(linesArr, arc)
		}
	}
	return linesArr
}

// Back fill all arcs to current time
func (r *RingClock) BackFillArcsContainer() {
	// startPointAngle := 0
	endPointAngleArr := r.timeAngles

	numBackFillArcs := []int{
		int(endPointAngleArr[0] / 6),  // seconds
		int(endPointAngleArr[1] / 6),  // minutes
		int(endPointAngleArr[2] / 30), // hours
	}

	angleStep := []float64{
		6,  // seconds step
		6,  // minutes step
		30, // hours step
	}

	for i, ring := range r.clockRings {
		cx := r.cx + r.offsetXarr[i]
		count := numBackFillArcs[i]
		stepDegrees := angleStep[i]

		for step := 0; step <= count; step++ {
			angle := float64(step) * stepDegrees
			x, y := getClockHandPosition(cx, r.cy, r.radius+r.arcOverShoot, angle)
			// arcs := r.drawArcsToResolution(cx, r.cy, x, y, r.arcStrokeWidth[i], ring.onColor)
			// for _, arc := range arcs {
			// 	r.arcContainers[i].Add(arc)
			// }

			arc := r.drawArc(cx, r.cy, x, y, r.arcStrokeWidth[i], ring.onColor)
			r.arcContainers[i].Add(arc)

		}
		r.refreshStaticClockItems(i)

	}

}

// Refreshes the clock static clock components that are over written by drawing the arcs
func (r *RingClock) refreshStaticClockItems(i int) {
	c := r.clocksContainer[i]

	cleaned := []fyne.CanvasObject{}
	for _, obj := range c.Objects {
		switch obj.(type) {
		case *canvas.Circle, *canvas.Text:
		default:
			cleaned = append(cleaned, obj)
		}
	}
	c.Objects = cleaned

	cx := r.cx + r.offsetXarr[i]

	innerCircle := canvas.NewCircle(theme.Color(theme.ColorNameBackground))
	innerCircle.Resize(fyne.NewSize(float32(r.radius*2-r.thickness), float32(r.radius*2-r.thickness)))
	innerCircle.Move(fyne.NewPos(
		float32(cx-r.radius+(r.thickness/2)),
		float32(r.cy-r.radius+(r.thickness/2)),
	))

	label := canvas.NewText(r.clockRings[i].Name, r.clockRings[i].labelColor)
	label.TextSize = float32(r.radius) * r.labelScaleFactorOfRadius
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Move(fyne.NewPos(
		float32(cx),
		float32(r.cy-(int(label.MinSize().Height)/2)),
	))

	// outer ring mask to cover overshoot from arcs to give perfect circle finish
	const maskStrokeWidth = float32(10)
	mask := canvas.NewCircle(color.Transparent)
	mask.StrokeWidth = maskStrokeWidth
	mask.StrokeColor = theme.Color(theme.ColorNameBackground)
	mask.Resize(fyne.NewSize(float32(r.radius*2+int(maskStrokeWidth)), float32(r.radius*2+int(maskStrokeWidth))))
	mask.Move(fyne.NewPos(float32(cx-r.radius-int(maskStrokeWidth/2)), float32(r.cy-r.radius-int(maskStrokeWidth/2))))

	r.clocksContainer[i].Add(mask)
	c.Add(innerCircle)
	c.Add(label)
	c.Add(mask)
	c.Refresh()

}

func (r *RingClock) Update(t *TickData) {
	angleArr := t.GetClockAngles()

	for i, ring := range r.clockRings {
		var arc *canvas.Line
		cx := r.cx + r.offsetXarr[i]
		x, y := getClockHandPosition(cx, r.cy, r.radius+r.arcOverShoot, angleArr[i])

		// Reset arc container based on ring type
		if ring.Name == "Seconds" && t.Second%60 == 0 {
			r.arcContainers[i].Objects = nil
			r.arcContainers[i].Refresh()
		}
		if ring.Name == "Minutes" && t.Minute%60 == 0 {
			r.arcContainers[i].Objects = nil
			r.arcContainers[i].Refresh()
		}
		if ring.Name == "Hours" && t.Hour12%12 == 0 {
			r.arcContainers[i].Objects = nil
			r.arcContainers[i].Refresh()
		}
		// Draw if tick detected
		if ring.Name == "Seconds" && t.SecondChanged() {
			arc = r.drawArc(cx, r.cy, x, y, r.arcStrokeWidth[i], ring.onColor)
		}
		if ring.Name == "Minutes" && t.MinuteChanged() {
			arc = r.drawArc(cx, r.cy, x, y, r.arcStrokeWidth[i], ring.onColor)
		}
		if ring.Name == "Hours" && t.HourChanged() {
			arc = r.drawArc(cx, r.cy, x, y, r.arcStrokeWidth[i], ring.onColor)
		}

		if arc != nil {
			r.arcContainers[i].Add(arc)
		}

		r.refreshStaticClockItems(i)

		// r.DebugPrintContainerCounts()
		r.DebugPrintArcContainers()
	}

}

func (r *RingClock) DebugPrintContainerCounts() {
	for i, container := range r.clocksContainer {
		fmt.Printf("[%s] has %d clock objects\n", r.clockRings[i].Name, len(container.Objects))
	}
	fmt.Println("..")
	for i, container := range r.arcContainers {
		fmt.Printf("[%s] has %d arc objects\n", r.clockRings[i].Name, len(container.Objects))
	}

	fmt.Println("----")
}

func (r *RingClock) DebugPrintArcContainers() {
	fmt.Println("Arc Containers Status:")
	for i, arcContainer := range r.arcContainers {
		fmt.Printf("  [%s] contains %d objects\n", r.clockRings[i].Name, len(arcContainer.Objects))
		for j, obj := range arcContainer.Objects {
			switch o := obj.(type) {
			case *canvas.Line:
				fmt.Printf("    Object %d: *canvas.Line StrokeWidth=%.2f Color=%v\n", j, o.StrokeWidth, o.StrokeColor)
			case *canvas.Circle:
				fmt.Printf("    Object %d: *canvas.Circle Color=%v\n", j, o.FillColor)
			case *canvas.Text:
				fmt.Printf("    Object %d: *canvas.Text Text=%q\n", j, o.Text)
			default:
				fmt.Printf("    Object %d: Unknown type %T\n", j, o)
			}
		}
	}
	fmt.Println("----")
}
