package main

import "os"
import "log"

// import "fmt"

import "github.com/gdamore/tcell"

func main() {

	var buf string = ""

	pos := map[string]interface{}{
		"x1": 7,
		"y1": 0,
	}

	TAG := map[string]interface{}{
		"label": "[dsh]$",
		"x1":    0,
		"y1":    0,
	}

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	s.Clear()

	drawText(s,
		TAG["x1"].(int),
		TAG["y1"].(int),
		defStyle,
		TAG["label"].(string))

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	for {
		// Update Screen
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == 127 {
				if pos["x1"].(int) > 7 {
					pos["x1"] = pos["x1"].(int) - 1

					buf = buf[:len(buf)-1]

					drawText(s,
						pos["x1"].(int),
						pos["y1"].(int),
						defStyle,
						" ",
					)
				}
			} else if ev.Key() == tcell.KeyEnter {
				if buf == "clear" {
					s.Clear()
					buf = ""

					TAG["y1"] = -1

					pos["y1"] = -1
					pos["x1"] = 7
				} else if buf == "exit" {
					quit()
				} else {
					buf = ""
				}
				// Update column position for dsh tag
				TAG["y1"] = TAG["y1"].(int) + 1

				// Update column and reset row position for cursor
				pos["y1"] = pos["y1"].(int) + 1
				pos["x1"] = 7
				drawText(s,
					TAG["x1"].(int),
					TAG["y1"].(int),
					defStyle,
					TAG["label"].(string),
				)
			} else if ev.Rune() >= 32 && ev.Rune() <= 126 {
				drawText(s,
					pos["x1"].(int),
					pos["y1"].(int),
					defStyle,
					string(rune(ev.Rune())),
				)
				pos["x1"] = pos["x1"].(int) + 1
				buf = buf + string(rune(ev.Rune()))
			}
		}
	}
}

func drawText(s tcell.Screen, x1, y1 int, style tcell.Style, text string) {
	row := y1
	col := x1

	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}
