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
	}
}
