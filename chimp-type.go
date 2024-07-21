package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const wordsPerLine = 10

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("there has been an error, oh oh oh... %v\n", err)
		os.Exit(1)
	}
}
