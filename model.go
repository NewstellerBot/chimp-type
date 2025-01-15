package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	MainMenu view = iota
	TypingTest
	Score
)

type model struct {
	// For storing text to be written
	prevLine     []string
	currentLine  []string
	nextLine     []string
	nextNextLine []string

	// For what user has entered
	prevUserInput []string
	userInput     []string

	// For measuring time
	startTime time.Time
	timeLeft  time.Duration

	// Keeping track of the view
	view view

	// Choices for tests
	choices []TestMode
	cursor  int

	// Score
	score int64
}

func initialModel() model {
	return model{
		prevLine:     nil,
		currentLine:  generateRandomFrom200(wordsPerLine),
		nextLine:     generateRandomFrom200(wordsPerLine),
		nextNextLine: nil,

		prevUserInput: nil,
		userInput:     []string{""},
		startTime:     time.Time{},

		view:    MainMenu,
		cursor:  0,
		choices: testModeChoices,
		score:   0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) UpdateTyping(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.choices[m.cursor] != zen {
		// Start time on the first input
		if m.startTime.IsZero() {
			m.startTime = time.Now()
		}

		// Update time left
		m.timeLeft = m.choices[m.cursor].Seconds() - time.Since(m.startTime)
		if m.choices[m.cursor] == zen {
			m.timeLeft = time.Hour
		}

		if m.timeLeft <= time.Microsecond {
			m.view = Score
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case tea.KeyBackspace.String():
			if len(m.userInput[len(m.userInput)-1]) > 0 {
				// Delete character within the word
				m.userInput[len(m.userInput)-1] = m.userInput[len(m.userInput)-1][:len(m.userInput[len(m.userInput)-1])-1]
			} else if len(m.userInput) > 1 {
				// Delete new word
				// if the last word is correct do not delete
				if m.IsLastWordCorrect() {
					return m, nil
				}
				m.userInput = m.userInput[:len(m.userInput)-1]
			} else if m.prevUserInput != nil {
				if m.IsLastWordCorrect() {
					return m, nil
				}
				// Move to previous line
				m.nextNextLine = m.nextLine
				m.nextLine = m.currentLine
				m.currentLine = m.prevLine
				m.prevLine = nil

				m.userInput = m.prevUserInput
				m.prevUserInput = nil
			}
		case tea.KeySpace.String():
			if len(m.userInput) == len(m.currentLine) {
				if m.prevLine != nil {
					m.score += int64(CalcLineScore(m.prevUserInput, m.prevLine))
				}
				m.prevLine = m.currentLine
				m.currentLine = m.nextLine

				m.prevUserInput = m.userInput
				m.userInput = []string{""}

				if m.nextNextLine != nil {
					m.nextLine = m.nextNextLine
					m.nextNextLine = nil
				} else {
					m.nextLine = generateRandomFrom200(wordsPerLine)
				}
			} else {
				m.userInput = append(m.userInput, "")
			}
		default:
			// If something that is not a letter is clicked, ignore it
			if len(msg.String()) > 1 {
				break
			}
			m.userInput[len(m.userInput)-1] += msg.String()
		}
	}

	return m, nil
}

func (m model) IsLastWordCorrect() bool {

	var lastCorr string
	var lastUser string

	if len(m.userInput) > 1 {
		lastCorr = m.currentLine[len(m.userInput)-2]
		lastUser = m.userInput[len(m.userInput)-2]
	} else if m.prevUserInput != nil {
		lastCorr = m.prevLine[len(m.prevUserInput)-1]
		lastUser = m.prevUserInput[len(m.prevUserInput)-1]
	} else {
		// First word
		return true
	}

	if len(lastCorr) != len(lastUser) {
		return false
	}

	_ = os.WriteFile("debug.txt", []byte(lastCorr+"\n\n"+lastUser), 0644)

	for i, ch := range lastUser {
		if lastCorr[i] != byte(ch) {
			return false
		}
	}

	return true
}

func (m model) StartTypingTest(t TestMode) {

}

func (m model) UpdateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k", "up":
			m.cursor = max(m.cursor-1, 0)
		case "j", "down":
			m.cursor = min(m.cursor+1, len(m.choices)-1)
		case "enter", " ":
			m.view = TypingTest
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.view {
	case MainMenu:
		return m.UpdateMainMenu(msg)
	case TypingTest:
		return m.UpdateTyping(msg)
	case Score:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) ViewTyping() string {
	s := "ChimpType\n\n"
	if m.choices[m.cursor] != zen {
		s += fmt.Sprintf("Time left: %.2f\n\n", float32(m.timeLeft.Microseconds())/1e6)
	}

	if m.prevUserInput != nil {
		for i, w := range m.prevUserInput {
			for j, ch := range w {
				if j < len(m.prevLine[i]) && byte(ch) == m.prevLine[i][j] {
					s += formatCorrect(ch)
				} else {
					if j < len(m.prevLine[i]) {
						s += formatIncorrect(m.prevLine[i][j])
					} else {
						s += formatExtraRune(ch)
					}
				}
			}

			if len(w) < len(m.prevLine[i]) {
				s += formatNotTyped(m.prevLine[i][len(w):])
			}
			s += " "
		}
	}

	if m.nextNextLine == nil {
		s += "\n\n"
	}

	displayString := ""

	for i, w := range m.userInput {
		for j, ch := range w {
			if j < len(m.currentLine[i]) && byte(ch) == m.currentLine[i][j] {
				displayString += formatCorrect(ch)
			} else {
				if j < len(m.currentLine[i]) {
					displayString += formatIncorrect(m.currentLine[i][j])
				} else {
					displayString += formatExtraRune(ch)
				}
			}
		}
		if len(w) < len(m.currentLine[i]) && len(m.userInput)-1 == i {
			// next char is formatted as cursor
			displayString += formatCursor(string(m.currentLine[i][len(w)]))
			// rest is normal
			if len(m.currentLine[i]) > len(w) {
				displayString += m.currentLine[i][len(w)+1:]
			}
		} else if len(w) < len(m.currentLine[i]) {
			displayString += formatNotTyped(m.currentLine[i][len(w):])
		}

		if i == len(m.userInput)-1 && len(m.currentLine[i]) <= len(w) {
			displayString += formatCursor(" ")
		} else {
			displayString += " "
		}
	}

	for i := len(m.userInput); i < len(m.currentLine); i++ {
		displayString += fmt.Sprintf("%v ", m.currentLine[i])
	}

	s += displayString + "\n\n"

	s += strings.Join(m.nextLine, " ")
	s += "\n\n"

	if m.nextNextLine != nil {
		s += strings.Join(m.nextNextLine, " ")
		s += "\n\n"
	}

	s += "\n\nPress ctrl+c to quit.\n"

	return s
}

func (m model) ViewMainMenu() string {
	s := "\nChimpType--minimal terminal typing test\n\nPlease select mode\n"

	for i, choice := range m.choices {

		cursor := " "
		if i == int(m.cursor) {
			cursor = ">"
		}

		s += cursor + " " + choice.String() + "\n"
	}
	s += "\n"
	return s
}

func (m model) CalculateFinalScore() float64 {
	s := m.score
	if m.prevLine != nil {
		s += int64(CalcLineScore(m.prevUserInput, m.prevLine))
	}
	s += int64(CalcLineScore(m.currentLine, m.currentLine))
	return float64(s) / CharsPerWord * (60.0 / m.choices[m.cursor].Seconds().Seconds())
}

func (m model) View() string {
	switch m.view {
	case MainMenu:
		return m.ViewMainMenu()
	case TypingTest:
		return m.ViewTyping()
	case Score:
		if m.choices[m.cursor] == zen {
			return "Zennnnnn...\n"
		}
		return fmt.Sprintf("Final score: %.2f wpm\n\nPlease click anything to continue...", m.CalculateFinalScore())
	}
	return "Nothing to see here...\n"
}
