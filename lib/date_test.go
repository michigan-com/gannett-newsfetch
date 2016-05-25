package lib

import (
	"testing"
	"time"
)

type TimeTestCase struct {
	TimeOne time.Time
	TimeTwo time.Time
}

func TestSameTime(t *testing.T) {
	now := time.Now()
	nextMillisecond := now.Add(time.Millisecond * 1)
	nextSecond := now.Add(time.Second * 1)
	nextMinute := now.Add(time.Minute * 1)
	nextHour := now.Add(time.Hour * 1)
	sameTimes := []TimeTestCase{
		TimeTestCase{now, now},
		TimeTestCase{now, nextMillisecond},
	}

	diffTimes := []TimeTestCase{
		TimeTestCase{now, nextSecond},
		TimeTestCase{now, nextMinute},
		TimeTestCase{now, nextHour},
	}

	dateTests(t, sameTimes, true)
	dateTests(t, diffTimes, false)
}

func dateTests(t *testing.T, times []TimeTestCase, expectedOutput bool) {

	for _, time := range times {
		if SameTime(time.TimeOne, time.TimeTwo) != expectedOutput {
			errorMsg := "Dates should match"
			if !expectedOutput {
				errorMsg = "Dates souldn't match"
			}
			t.Fatalf("%s: %v & %v", errorMsg, time.TimeOne, time.TimeTwo)
		}
	}
}
