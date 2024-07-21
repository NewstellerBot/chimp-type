package main

import (
	"fmt"
	"testing"
)

func TestCalcLineScore(t *testing.T) {
	var tests = []struct {
		c, u []string
		want int
	}{
		{[]string{"1234", "1234"}, []string{"1234", "1234"}, 10},
		{[]string{"1234", "1234"}, []string{"12", "12"}, 0},
		{[]string{"1234", "1234"}, []string{"1abc", "1abc"}, 0},
		{[]string{"1234", "1234"}, []string{"1234", "123"}, 5},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v,%v", tt.c, tt.u)
		t.Run(testname, func(t *testing.T) {
			ans := CalcLineScore(tt.c, tt.u)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
