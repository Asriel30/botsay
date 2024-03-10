package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/mattes/go-asciibot"
	"github.com/xyproto/rainbow"
)

const (
	boxContentWidth = 42
	versionString   = "botsay 1.2.7"
)

// GFX is ASCII graphics as a string, and where to place it on the canvas
type GFX struct {
	ascii string
	x     int
	y     int
}

// New creates a new GFX struct, with an ASCII art string and a position
func New(ascii string, x, y int) *GFX {
	return &GFX{ascii, x, y}
}

// Draw an ASCII bubble
func bubble(w, h int) string {
	var (
		sb     strings.Builder
		dashes = strings.Repeat("-", w-5)
	)
	sb.WriteString("   .")
	sb.WriteString(dashes)
	sb.WriteString(".\n")
	for i := 0; i < (h - 2); i++ {
		if i == 1 {
			sb.WriteString("--<|")
		} else {
			sb.WriteString("   |")
		}
		sb.WriteString(strings.Repeat(" ", w-5))
		sb.WriteString("|\n")
	}
	sb.WriteString("   '")
	sb.WriteString(dashes)
	sb.WriteString("'\n")
	return sb.String()
}

// Combine several ASCII graphics layers (with a position each) into one layer
func render(layers []*GFX) string {
	var canvas string
	for _, gfx := range layers {
		canvas = CombineArt(canvas, gfx.ascii, gfx.x, gfx.y)
	}
	return canvas
}

// Generate ASCII graphics of a randomly generated bot with a speech bubble
func botsay(msg string) string {
	var layers []*GFX
	trimmed := strings.TrimSpace(msg)
	msgwidth := boxContentWidth
	lineCount := strings.Count(trimmed, "\n") + 1
	layers = append(layers, New(asciibot.Random(), 1, 1))
	sl := SplitWidthWords(trimmed, msgwidth)
	boxX := 18
	boxY := 1
	if RuneLen(trimmed) > 0 {
		layers = append(layers, New(bubble(min(msgwidth, RuneLen(trimmed))+7, len(sl)+lineCount+1), boxX, boxY))
		counter := 0
		for _, s := range sl {
			layers = append(layers, New(s, boxX+5, boxY+1+counter))
			counter++
		}
	}
	return strings.TrimRightFunc(render(layers), unicode.IsSpace) + "\n"
}

func main() {
	rainbowMode := false
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
		}
		switch args[0] {
		case "--help":
			fmt.Println("usage: botsay [-c] [TEXT or \"-\"]")
			return
		case "--version":
			fmt.Println(versionString)
			return
		case "-c":
			rainbowMode = true
			if len(args) > 1 {
				args = args[1:]
				if len(args) > 0 && args[0] == "--" {
					args = args[1:]
				}
			} else {
				args = []string{}
			}
		}
	}
	// Join all arguments to a single string
	msg := strings.Join(args, " ")
	// Read from /dev/stdin if "-" is given
	if msg == "-" {
		data, err := os.ReadFile("/dev/stdin")
		if err != nil {
			panic(err)
		}
		msg = strings.TrimSpace(string(data))
	}
	if rainbowMode {
		rw := rainbow.NewTruecolorWriter(3, 0.4, 10)
		rw.Write([]byte(botsay(msg) + "\n"))
	} else {
		fmt.Println(botsay(msg))
	}
}
