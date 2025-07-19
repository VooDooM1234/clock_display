package clock

import (
	"time"
)

func GetTimeConstruct() (currenTime time.Time, hh int, mm int, ss int) {
	currentTime := time.Now()

	hh = currentTime.Hour()
	mm = currentTime.Minute()
	ss = currentTime.Second()

	return currentTime, hh, mm, ss
}

func drawHand(cx int, cy int, x int, y int, rad int) {

}

func drawClock(cx int, cy int, radius int) {

}
