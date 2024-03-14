package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/mattes/go-asciibot"
	"github.com/spf13/pflag"
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

// bubble will draw an ASCII bubble
func bubble(w, h int) string {
	var (
		sb     strings.Builder
		dashes = strings.Repeat("-", w-5)
	)
	sb.WriteString("   ." + dashes + ".\n")
	for i := 0; i < (h - 2); i++ {
		if i == 1 {
			sb.WriteString("--<|")
		} else {
			sb.WriteString("   |")
		}
		sb.WriteString(strings.Repeat(" ", w-5) + "|\n")
	}
	sb.WriteString("   '" + dashes + "'\n")
	return sb.String()
}

// render will combine several ASCII graphics layers (with a position each) into a single layer
func render(layers []*GFX) string {
	var canvas []string
	// Determine the required dimensions of the canvas
	maxWidth, maxHeight := 0, 0
	for _, gfx := range layers {
		gfxWidth, gfxHeight := Dimensions(gfx.ascii)
		if gfx.x+gfxWidth > maxWidth {
			maxWidth = gfx.x + gfxWidth
		}
		if gfx.y+gfxHeight > maxHeight {
			maxHeight = gfx.y + gfxHeight
		}
	}

	// Initialize the canvas with spaces
	canvas = make([]string, maxHeight)
	for i := range canvas {
		canvas[i] = strings.Repeat(" ", maxWidth)
	}

	for _, gfx := range layers {
		gfxLines := strings.Split(gfx.ascii, "\n")
		for y, line := range gfxLines {
			canvasY := gfx.y + y
			if canvasY >= len(canvas) {
				continue // Skip if this line is out of the canvas bounds
			}
			for x, ch := range line {
				canvasX := gfx.x + x
				if canvasX >= len(canvas[canvasY]) {
					continue // Skip if this character is out of the line bounds
				}
				// Replace the character in the canvas
				canvasLine := []rune(canvas[canvasY])
				canvasLine[canvasX] = ch
				canvas[canvasY] = string(canvasLine)
			}
		}
	}
	return strings.Join(canvas, "\n")
}

// botsay will generate ASCII graphics of the specified bot ID, and with a speech bubble
func botsay(msg string, botID string) string {
	var layers []*GFX
	trimmed := strings.TrimSpace(msg)
	msgWidth := boxContentWidth
	lineCount := strings.Count(trimmed, "\n") + 1
	botASCII, _ := asciibot.Generate(botID)
	layers = append(layers, New(botASCII, 1, 1))
	sl := SplitWidthWords(trimmed, msgWidth)
	boxX := 18
	boxY := 1
	if RuneLen(trimmed) > 0 {
		layers = append(layers, New(bubble(min(msgWidth, RuneLen(trimmed))+7, len(sl)+lineCount+1), boxX, boxY))
		for counter, s := range sl {
			layers = append(layers, New(s, boxX+5, boxY+1+counter))
		}
	}
	return strings.TrimRightFunc(render(layers), unicode.IsSpace) + "\n"
}

func main() {
	var (
		rainbowMode bool
		customBotID string
		printID     bool
		versionFlag bool
		helpFlag    bool
	)

	pflag.BoolVarP(&rainbowMode, "color", "c", false, "Enable rainbow mode")
	pflag.StringVarP(&customBotID, "id", "i", "", "Specify a custom bot ID to use for generating the ASCII art.")
	pflag.BoolVarP(&printID, "print", "p", false, "Print the bot's ID after generating the ASCII art.")
	pflag.BoolVar(&versionFlag, "version", false, "Print the version and exit")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "Show this help message")
	pflag.Parse()

	if versionFlag {
		fmt.Println(versionString)
		return
	}

	if helpFlag {
		pflag.Usage()
		return
	}

	botID := customBotID
	if botID == "" {
		botID = asciibot.RandomID()
	}

	args := pflag.Args()
	msg := strings.Join(args, " ")
	output := botsay(msg, botID)

	if rainbowMode {
		rw := rainbow.NewTruecolorWriter(3, 0.4, 10)
		rw.Write([]byte(output + "\n"))
	} else {
		fmt.Println(output)
	}

	if printID {
		fmt.Println("Bot ID:", botID)
	}
}
