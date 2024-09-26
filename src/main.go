package main

import "os"
import "log"

// import "fmt"

import "github.com/gdamore/tcell"

func main() {

	var buf string = ""

	pos := map[string]interface{}{
		"x1": 7,
		"x2": 8,
		"y1": 0,
		"y2": 1,
	}

	TAG := map[string]interface{}{
		"label": "[dsh]$",
		"x1":    0,
		"x2":    6,
		"y1":    0,
		"y2":    1,
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
		TAG["x2"].(int),
		TAG["y2"].(int),
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
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyEnter {
				if buf == "clear" {
					s.Clear()
					buf = ""

					TAG["y1"] = -1
					TAG["y2"] = 0

					pos["y1"] = -1
					pos["y2"] = 1
					pos["x1"] = 7
					pos["x2"] = 8
				} else {
					buf = ""
				}
				// Update column position for dsh tag
				TAG["y1"] = TAG["y1"].(int) + 1
				TAG["y2"] = TAG["y2"].(int) + 1

				// Update column and reset row position for cursor
				pos["y1"] = pos["y1"].(int) + 1
				pos["y2"] = pos["y2"].(int) + 1
				pos["x1"] = 7
				pos["x2"] = 8
				drawText(s,
					TAG["x1"].(int),
					TAG["y1"].(int),
					TAG["x2"].(int),
					TAG["x2"].(int),
					defStyle,
					TAG["label"].(string),
				)
			} else {
				drawText(s,
					pos["x1"].(int),
					pos["y1"].(int),
					pos["x2"].(int),
					pos["y2"].(int),
					defStyle,
					string(rune(ev.Rune())),
				)
				pos["x1"] = pos["x1"].(int) + 1
				pos["x2"] = pos["x2"].(int) + 1
				buf = buf + string(rune(ev.Rune()))
			}
		}
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1

	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row >= y2 {
			break
		}
	}
}
