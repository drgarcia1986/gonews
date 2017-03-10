package gui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	if name := v.Name(); name != "main" {
		g.Cursor = false
		return deleteView(name, g)
	}
	return gocui.ErrQuit
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	cx, cy := v.Cursor()
	if line, err := v.Line(cy + 1); err != nil || line == "" {
		return nil
	}

	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}

	return nil
}

func helpMsg(g *gocui.Gui, _ *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("help", maxX/2-30, maxY/2, maxX/2+15, maxY/2+6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keyboard Shortcuts"
		fmt.Fprintln(v, `
		- Enter : Open story in default browser
		- p : Preview (beta)
		- k / Arrow Up : Move up
		- j / Arrow Down : Move down
		- q / Ctrl+c : Close window
		- ? : Show this message
		`)

		if _, err := g.SetCurrentView("help"); err != nil {
			return err
		}
	}
	return nil
}

func deleteView(viewName string, g *gocui.Gui) error {
	if err := g.DeleteView(viewName); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}
	return nil
}

func showPreview(g *gocui.Gui, title, content string) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("preview", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprint("Preview - ", title)
		g.Cursor = true

		var length int
		for _, s := range strings.Split(content, "\n") {
			for {
				length = maxX - 2
				if length > len(s) {
					length = len(s)
				}

				text := strings.TrimSpace(s[:length])
				if text != "" {
					fmt.Fprintln(v, text)
				}

				if length == len(s) {
					break
				}
				s = s[length:]
			}
		}

		if _, err := g.SetCurrentView("preview"); err != nil {
			return err
		}
	}
	return nil
}
