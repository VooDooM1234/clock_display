// Tick Data hold
package clock

import (
	"fmt"
	"strconv"
	"time"
)

type TickData struct {
	Hour12        int
	Hour24        int
	Minute        int
	Second        int
	Hr12OnesDigit int
	Hr12TensDigit int
	Hr24OnesDigit int
	Hr24TensDigit int
	MinOnesDigit  int
	MinTensDigit  int
	SecOnesDigit  int
	SecTensDigit  int
	prevSec       int
	prevMin       int
	prevHr        int
}

func NewTickData() *TickData {
	now := time.Now()

	Hour12 := (now.Hour()) % 12
	Hour24 := now.Hour()
	Minute := now.Minute()
	Second := now.Second()

	hour12Str := fmt.Sprintf("%02d", Hour12)
	hour24Str := fmt.Sprintf("%02d", Hour24)
	minuteStr := fmt.Sprintf("%02d", Minute)
	secondStr := fmt.Sprintf("%02d", Second)

	Hr12TensDigit, _ := strconv.Atoi(string(hour12Str[0]))
	Hr12OnesDigit, _ := strconv.Atoi(string(hour12Str[1]))

	Hr24TensDigit, _ := strconv.Atoi(string(hour24Str[0]))
	Hr24OnesDigit, _ := strconv.Atoi(string(hour24Str[1]))

	MinTensDigit, _ := strconv.Atoi(string(minuteStr[0]))
	MinOnesDigit, _ := strconv.Atoi(string(minuteStr[1]))

	SecTensDigit, _ := strconv.Atoi(string(secondStr[0]))
	SecOnesDigit, _ := strconv.Atoi(string(secondStr[1]))

	return &TickData{
		Hour12:        Hour12,
		Hour24:        Hour24,
		Minute:        Minute,
		Second:        Second,
		Hr12OnesDigit: Hr12OnesDigit,
		Hr12TensDigit: Hr12TensDigit,
		Hr24OnesDigit: Hr24OnesDigit,
		Hr24TensDigit: Hr24TensDigit,

		MinOnesDigit: MinOnesDigit,
		MinTensDigit: MinTensDigit,
		SecOnesDigit: SecOnesDigit,
		SecTensDigit: SecTensDigit,

		prevSec: -1,
		prevMin: -1,
		prevHr:  -1,
	}
}

func (time *TickData) GetClockAngles() []float64 {

	var anglesArr []float64

	hourAngle := (float64(time.Hour12) + float64(time.Minute)/60.0) * 30.0
	minuteAngle := float64(time.Minute) * 6
	secondAngle := float64(time.Second) * 6

	anglesArr = append(anglesArr, secondAngle, minuteAngle, hourAngle)

	return anglesArr
}

func (t *TickData) Update() {
	now := time.Now()

	// Save previous values
	t.prevHr = t.Hour12
	t.prevMin = t.Minute
	t.prevSec = t.Second

	// Update current values
	t.Hour12 = now.Hour() % 12
	t.Hour24 = now.Hour()
	t.Minute = now.Minute()
	t.Second = now.Second()

	// Update digit fields (optional, keep if needed)
	hour12Str := fmt.Sprintf("%02d", t.Hour12)
	hour24Str := fmt.Sprintf("%02d", t.Hour24)
	minuteStr := fmt.Sprintf("%02d", t.Minute)
	secondStr := fmt.Sprintf("%02d", t.Second)

	t.Hr12TensDigit, _ = strconv.Atoi(string(hour12Str[0]))
	t.Hr12OnesDigit, _ = strconv.Atoi(string(hour12Str[1]))

	t.Hr24TensDigit, _ = strconv.Atoi(string(hour24Str[0]))
	t.Hr24OnesDigit, _ = strconv.Atoi(string(hour24Str[1]))

	t.MinTensDigit, _ = strconv.Atoi(string(minuteStr[0]))
	t.MinOnesDigit, _ = strconv.Atoi(string(minuteStr[1]))

	t.SecTensDigit, _ = strconv.Atoi(string(secondStr[0]))
	t.SecOnesDigit, _ = strconv.Atoi(string(secondStr[1]))
}

func (t *TickData) SecondChanged() bool {
	if t.prevSec == -1 {
		return true // First run, consider it changed
	}
	return t.Second != t.prevSec
}

func (t *TickData) MinuteChanged() bool {
	if t.prevMin == -1 {
		return true
	}
	return t.Minute != t.prevMin
}

func (t *TickData) HourChanged() bool {
	if t.prevHr == -1 {
		return true
	}
	return t.Hour12 != t.prevHr
}

// func GetClockAngles() []float64 {

// 	var anglesArr []float64

// 	hourAngle := (float64(time.Hour12) + float64(time.Minute)/60.0) * 30.0
// 	minuteAngle := float64(time.Minute) * 6
// 	secondAngle := float64(time.Second) * 6

// 	anglesArr = append(anglesArr, secondAngle, minuteAngle, hourAngle)

// 	return anglesArr
// }
