// Package svgscreen implements a fixed font terminal screen using SVG
package svgscreen

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/wader/ansisvg/color"
)

//go:embed template.svg
var templateSVG string

type Char struct {
	Char       string
	X          int
	Foreground string
	Background string
	Underline  bool
	Intensity  bool
	Invert     bool
}

type Line struct {
	Y     int
	Chars []Char
}

type BoxSize struct {
	Width  int
	Height int
}

type TextSpan struct {
	ForegroundColor string
	Content         string
}

type TextElement struct {
	Y         int
	TextSpans []TextSpan
}

type Screen struct {
	Transparent      bool
	ForegroundColor  string
	ForegroundColors map[string]string
	BackgroundColor  string
	BackgroundColors map[string]string
	FontName         string
	FontEmbedded     []byte
	FontRef          string
	FontSize         int
	CharacterBoxSize BoxSize
	TerminalWidth    int
	Columns          int
	NrLines          int
	Lines            []Line
	TextElements     []TextElement
}

// Convert a line into a <text> element
// fc gives (color, content) of a char
func LineToTextElement(s Screen, l Line, fc func(Char) (string, string)) TextElement {
	result := TextElement{
		Y: l.Y * s.CharacterBoxSize.Height,
	}
	currentColor := ""
	currentContent := ""
	appendSpan := func() {
		if currentContent == "" {
			return
		}
		result.TextSpans = append(result.TextSpans, TextSpan{
			ForegroundColor: currentColor,
			Content:         currentContent,
		})
		currentContent = ""
	}
	for _, c := range l.Chars {
		charColor, charContent := fc(c)
		if charColor != currentColor {
			appendSpan()
		}
		currentColor = charColor
		currentContent += charContent
	}
	appendSpan()
	return result
}

func Render(w io.Writer, s Screen) error {
	t := template.New("")
	t.Funcs(template.FuncMap{
		"add":          func(a int, b int) int { return a + b },
		"mul":          func(a int, b int) int { return a * b },
		"hasprefix":    strings.HasPrefix,
		"iswhitespace": func(a string) bool { return strings.TrimSpace(a) == "" },
		"coloradd": func(a string, b string) string {
			return color.NewFromHex(a).Add(color.NewFromHex(b)).Hex()
		},
		"base64": func(bs []byte) string { return base64.RawStdEncoding.EncodeToString(bs) },
	})

	for _, l := range s.Lines {
		for i, c := range l.Chars {
			if c.Invert {
				c.Background, c.Foreground = c.Foreground, c.Background
				if c.Background == "" {
					c.Background = s.ForegroundColor
				}
				if c.Foreground == "" {
					c.Foreground = s.BackgroundColor
				}
				l.Chars[i] = c
			}
		}
		s.TextElements = append(s.TextElements, LineToTextElement(s, l, func(c Char) (string, string) {
			if c.Background == "" {
				return "", "&nbsp;"
			} else {
				return c.Background, "&#x2588;"
			}
		}))
		s.TextElements = append(s.TextElements, LineToTextElement(s, l, func(c Char) (string, string) {
			return c.Foreground, c.Char
		}))
	}
	for _, t := range s.TextElements {
		fmt.Fprintf(os.Stderr, "text at %d\n", t.Y)
		for _, s := range t.TextSpans {
			fmt.Fprintf(os.Stderr, "tspan col:%s cont: %s\n", s.ForegroundColor, s.Content)
		}
	}

	t, err := t.Parse(templateSVG)
	if err != nil {
		return err
	}
	if err = t.ExecuteTemplate(w, "", s); err != nil {
		return err
	}

	return nil
}
