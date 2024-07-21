package main

// Calculates the number of correct characters in a line
// Spaces are added if the lengths of the words are
// the same, regardless of their correctness.
// Params
// c: correct line
// u: user generated line
func CalcLineScore(c, u []string) int {
	score := 0
	for i, w := range u {
		// If words are not the same length, they cannot be correct
		if len(w) == len(c[i]) {
			correct := true
			// Iterate over characters in a word
			for j, ch := range w {
				// If one of them is not correct
				if ch != rune(c[i][j]) || j >= len(c[i]) {
					correct = false
					break
				}
			}
			if correct {
				// Add to score for a valid space
				score++
				score += len(w)
			}
		}
	}
	return score
}

const CharsPerWord = 5.0
