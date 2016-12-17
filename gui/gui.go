package gui

import (
	"fmt"

	"github.com/drgarcia1986/gonews/hackernews"
	"github.com/drgarcia1986/gonews/utils"
	"github.com/jroimartin/gocui"
)

type Gui struct {
	items []*hackernews.Story
}

func (gui *Gui) getLine(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return nil
	}

	for _, story := range gui.items {
		if story.Title == line {
			return utils.OpenUrl(story.Url)
		}
	}
	return nil
}

func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Go - Hacker News"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		for _, story := range gui.items {
			fmt.Fprintln(v, story.Title)
		}

		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, gui.getLine); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(gui.layout)
	if err := gui.keybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func New(items []*hackernews.Story) *Gui {
	guiItems := make([]*hackernews.Story, len(items))
	copy(guiItems, items)
	return &Gui{items: guiItems}
}
