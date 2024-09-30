package main

import "os"
import "log"

// import "fmt"

import "github.com/gdamore/tcell"

type CursorPos struct {
	x int
	y int
}

func main() {

	var buf string = ""

	pos := CursorPos{
		x: 7,
		y: 0,
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

	drawText(s, 0,
		pos.y,
		defStyle,
		"[dsh]$",
	)
	s.ShowCursor(pos.x, pos.y)

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
			drawText(s,
				pos.x,
				pos.y,
				defStyle,
				" ",
			)
			if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				if pos.x > 7 {
					charPos := pos.x - 7
					rmChar(s, &pos, charPos, &buf)
				}
			} else if ev.Key() == tcell.KeyEnter {
				if buf == "clear" {
					s.Clear()
					buf = ""

					pos.y = -1
					pos.x = 7
				} else if buf == "exit" {
					quit()
				} else {
					if buf != "" {
						drawText(s,
							0,
							pos.y+1,
							defStyle,
							"dsh: command not found: "+buf,
						)

						pos.y++
						pos.x = 7
					}
					buf = ""
				}
				pos.y++

				drawText(s, 0,
					pos.y,
					defStyle,
					"[dsh]$",
				)
			} else if ev.Key() == tcell.KeyRight {
				if pos.x-7 < len(buf) {
					pos.x++
				}
			} else if ev.Key() == tcell.KeyLeft {
				if pos.x > 7 {
					pos.x--
				}
			} else if ev.Rune() >= 32 && ev.Rune() <= 126 {
				drawText(s,
					pos.x,
					pos.y,
					defStyle,
					string(rune(ev.Rune())),
				)
				pos.x++
				buf = buf + string(rune(ev.Rune()))
			}
			drawCursor(s, pos, buf)
		}
	}
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	row := y
	col := x

	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}

func drawCursor(s tcell.Screen, pos CursorPos, buf string) {
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	drawText(s,
		7,
		pos.y,
		defStyle,
		buf,
	)
	s.ShowCursor(pos.x, pos.y)
}

func rmChar(s tcell.Screen, pos *CursorPos, index int, buf *string) {
	if len(*buf) == 0 || pos.x <= 7 {
		return
	}

	*buf = (*buf)[:index-1] + (*buf)[index:]

	width, _ := s.Size()
	for x := 7; x < width; x++ {
		s.SetContent(x, pos.y, ' ', nil, tcell.StyleDefault)
	}

	drawText(s,
		7,
		pos.y,
		tcell.StyleDefault,
		*buf,
	)

	pos.x--
	if pos.x < 7 {
		pos.x = 7
	}
}
