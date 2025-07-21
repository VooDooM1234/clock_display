package clock

import (
	"fmt"
	"strconv"
	"time"
)

// TickData holds the current time components.
type TickData struct {
	Hour12       int
	Hour24       int
	Minute       int
	Second       int
	HrOnesDigit  int
	HrTensDigit  int
	MinOnesDigit int
	MinTensDigit int
	SecOnesDigit int
	SecTensDigit int
}

func NewTickData() *TickData {
	now := time.Now()

	Hour12 := (now.Hour()) % 12
	Hour24 := now.Hour()
	Minute := now.Minute()
	Second := now.Second()

	hourStr := fmt.Sprintf("%02d", Hour12)
	minuteStr := fmt.Sprintf("%02d", Minute)
	secondStr := fmt.Sprintf("%02d", Second)

	HrTensDigit, _ := strconv.Atoi(string(hourStr[0]))
	HrOnesDigit, _ := strconv.Atoi(string(hourStr[1]))

	MinTensDigit, _ := strconv.Atoi(string(minuteStr[0]))
	MinOnesDigit, _ := strconv.Atoi(string(minuteStr[1]))

	SecTensDigit, _ := strconv.Atoi(string(secondStr[0]))
	SecOnesDigit, _ := strconv.Atoi(string(secondStr[1]))

	return &TickData{
		Hour12:       Hour12,
		Hour24:       Hour24,
		Minute:       Minute,
		Second:       Second,
		HrOnesDigit:  HrOnesDigit,
		HrTensDigit:  HrTensDigit,
		MinOnesDigit: MinOnesDigit,
		MinTensDigit: MinTensDigit,
		SecOnesDigit: SecOnesDigit,
		SecTensDigit: SecTensDigit,
	}
}

func (t *TickData) Update() {
	now := time.Now()
	t.Hour12 = now.Hour() % 12
	t.Hour24 = now.Hour()
	t.Minute = now.Minute()
	t.Second = now.Second()

	hourStr := fmt.Sprintf("%02d", t.Hour12)
	minuteStr := fmt.Sprintf("%02d", t.Minute)
	secondStr := fmt.Sprintf("%02d", t.Second)

	t.HrTensDigit, _ = strconv.Atoi(string(hourStr[0]))
	t.HrOnesDigit, _ = strconv.Atoi(string(hourStr[1]))
	t.MinTensDigit, _ = strconv.Atoi(string(minuteStr[0]))
	t.MinOnesDigit, _ = strconv.Atoi(string(minuteStr[1]))
	t.SecTensDigit, _ = strconv.Atoi(string(secondStr[0]))
	t.SecOnesDigit, _ = strconv.Atoi(string(secondStr[1]))
}

// func setTime(t time.Time) *TickData {

// }
