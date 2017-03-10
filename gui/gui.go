package gui

import (
	"fmt"

	"github.com/drgarcia1986/gonews/story"
	"github.com/drgarcia1986/gonews/utils"
	"github.com/jroimartin/gocui"
)

var (
	downKeys = []interface{}{'j', gocui.KeyArrowDown}
	upKeys   = []interface{}{'k', gocui.KeyArrowUp}
	quitKeys = []interface{}{'q', gocui.KeyCtrlC}

	keybindingMap = []struct {
		keys     []interface{}
		viewName string
		event    func(*gocui.Gui, *gocui.View) error
	}{
		{quitKeys, "", quit},
		{downKeys, "", cursorDown},
		{upKeys, "", cursorUp},
		{[]interface{}{'?'}, "main", helpMsg},
	}
)

type Gui struct {
	items        []*story.Story
	providerName string
}

func (gui *Gui) openStory(g *gocui.Gui, v *gocui.View) error {
	s, err := gui.getStoryOfCurrentLine(v)
	if err == nil && s != nil {
		return utils.OpenURL(s.URL)
	}
	return err
}

func (gui *Gui) preview(g *gocui.Gui, v *gocui.View) error {
	s, err := gui.getStoryOfCurrentLine(v)
	if err != nil {
		return err
	}

	if s == nil {
		return nil
	}

	content, err := utils.GetPreview(s.URL)
	if err != nil {
		return err
	}

	if content == "" {
		content = "No preview available"
	}
	return showPreview(g, s.Title, content)
}

func (gui *Gui) getStoryOfCurrentLine(v *gocui.View) (*story.Story, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return nil, err
	}

	for _, s := range gui.items {
		if s.Title == line {
			return s, nil
		}
	}
	return nil, nil
}

func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("GoNews - %s ('?' for help)", gui.providerName)
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
	for _, bm := range keybindingMap {
		for _, key := range bm.keys {
			if err := g.SetKeybinding(bm.viewName, key, gocui.ModNone, bm.event); err != nil {
				return err
			}
		}
	}

	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, gui.openStory); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", 'p', gocui.ModNone, gui.preview); err != nil {
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

	g.SetManagerFunc(gui.layout)
	if err := gui.keybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func New(items []*story.Story, providerName string) *Gui {
	guiItems := make([]*story.Story, len(items))
	copy(guiItems, items)
	return &Gui{items: guiItems, providerName: providerName}
}
