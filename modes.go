package main

import (
	"time"
)

type TestMode int

const (
	S15 = iota + 1
	S30
	S45
	S60
	zen
)

var testModeChoices = []TestMode{S15, S30, S45, S60, zen}

var modeName = map[TestMode]string{
	zen: "zen",
	S15: "15 seconds",
	S30: "30 seconds",
	S45: "45 seconds",
	S60: "60 seconds",
}

var modeSeconds = map[TestMode]time.Duration{
	zen: -1 * time.Second,
	S15: 15 * time.Second,
	S30: 30 * time.Second,
	S45: 45 * time.Second,
	S60: 60 * time.Second,
}

func (m TestMode) Seconds() time.Duration {
	return modeSeconds[m]
}

func (m TestMode) String() string {
	return modeName[m]
}
