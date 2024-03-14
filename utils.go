package main

import (
	"strings"
	"unicode/utf8"
)

// RuneLen returns the rune count, not the byte count
func RuneLen(s string) int {
	return utf8.RuneCountInString(s)
}

// Dimensions returns the width and height of a given ASCII art string
func Dimensions(asciiArt string) (int, int) {
	maxWidth := 0
	maxHeight := 0
	lineCounter := 0
	for _, line := range strings.Split(asciiArt, "\n") {
		l := RuneLen(line)
		if l > maxWidth {
			maxWidth = l
		}
		lineCounter++
	}
	if lineCounter > maxHeight {
		maxHeight = lineCounter
	}
	return maxWidth, maxHeight
}

// get returns a character at (x,y) in a multiline string.
// If anything go wrong, or if (x,y) is out of bounds, it return a space rune.
func get(s []rune, x, y, w, h int) rune {
	if x < 0 || y < 0 {
		return ' '
	}
	if x >= w || y >= h {
		return ' '
	}
	// +1 to account for the trailing newlines
	pos := y*w + x
	if pos >= len(s) {
		return ' '
	}
	r := s[pos]
	switch r {
	case '\n', '\t', '\r', '\v':
		return ' '
	}
	return r
}

// toMap can convert from a multiline-string to an indexed slice of runes (y*w+x style)
func toMap(s string, w int) []rune {
	rs := make([]rune, 0)
	for _, line := range strings.Split(s, "\n") {
		rs = append(rs, []rune(line)...)
		linelen := RuneLen(line)
		if linelen < w {
			// Fill out the rest of the line with spaces
			rs = append(rs, []rune(strings.Repeat(" ", w-linelen))...)
		}
	}
	return rs
}

// SplitWords can split a string into words, keeping punctuation and trailing spaces
func SplitWords(s string) []string {
	var (
		splitpoint bool
		words      []string
		letters    strings.Builder
		tmp        string
	)
	lenS := RuneLen(s)
	for i, r := range s {
		splitpoint = false
		switch r {
		case '.', '!', ',', ':', '-', ' ', '?', ';', '\n':
			// Check if the next character is not an end quote
			if i+1 < lenS && s[i+1] != '"' && s[i+1] != '\'' {
				splitpoint = true
			}
		}
		// Combine repeated dashes
		if r == '-' && i+1 < lenS && s[i+1] == '-' {
			splitpoint = false
		}
		// Combine repeated dots
		if r == '.' && i+1 < lenS && s[i+1] == '.' {
			splitpoint = false
		}
		if splitpoint || i == lenS {
			letters.WriteRune(r)
			tmp = letters.String()
			if RuneLen(tmp) > 0 {
				words = append(words, tmp)
			}
			letters.Reset()
		} else {
			letters.WriteRune(r)
		}
	}
	tmp = strings.TrimSpace(letters.String())
	if RuneLen(tmp) > 0 {
		words = append(words, tmp)
	}
	return words
}

// SplitWidthWords can split a string by words, then combine to form lines maximum w long
func SplitWidthWords(s string, w int) []string {
	var (
		sl   []string
		line string
	)
	for _, word := range SplitWords(s) {
		if RuneLen(line)+RuneLen(word) < w {
			line += word
		} else {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasSuffix(trimmedLine, "--") {
				// Move the double dash to the beginning of the next line
				trimmedLine = trimmedLine[:RuneLen(trimmedLine)-2]
				sl = append(sl, trimmedLine)
				line = "-- " + word
			} else {
				sl = append(sl, trimmedLine)
				line = word
			}
		}
	}
	if RuneLen(line) == 0 {
		return sl
	}
	return append(sl, strings.TrimSpace(line))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
