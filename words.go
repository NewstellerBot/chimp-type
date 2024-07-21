package main

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed top200.txt
var top200 string

func generateRandomFrom200(n int) []string {
	res := []string{}
	top200list := strings.Split(top200, "\n")
	for range n {
		randomIx := rand.Intn(200)
		res = append(res, top200list[randomIx])
	}

	return res
}
