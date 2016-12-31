package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == "help" {
		return closeHelpMsg(g, v)
	}
	return gocui.ErrQuit
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	cx, cy := v.Cursor()
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

func helpMsg(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("help", maxX/2-30, maxY/2, maxX/2+15, maxY/2+6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keyboard Shortcuts"
		fmt.Fprintln(v, `
		- Enter : Open story in default browser
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

func closeHelpMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("help"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}
	return nil
}
