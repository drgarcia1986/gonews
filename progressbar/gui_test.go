package progressbar

import (
	"strings"
	"testing"
)

func TestGUIStartNew(t *testing.T) {
	p := new(GUIProgressBar)

	expected := 10
	p.Start(expected)

	if p.p == nil {
		t.Fatal("error on start new progressbar")
	}

	if p.p.Total != int64(expected) {
		t.Errorf("expected %d, got %d", expected, p.p.Total)
	}
}

func TestGUIIncrement(t *testing.T) {
	p := new(GUIProgressBar)

	p.Start(10)
	p.Increment()

	if p.p.Get() != int64(1) {
		t.Errorf("expected 1, got %d", p.p.Get())
	}
}

func TestGUIFinish(t *testing.T) {
	p := New()

	p.Start(1)
	p.Increment()
	p.Finish()

	g, ok := p.(*GUIProgressBar)
	if !ok {
		t.Fatal("error in type assertation of progress interface")
	}

	expected := "100.00%"
	if !strings.Contains(g.p.String(), "100.00%") {
		t.Errorf("expected %s, got %s", expected, g.p.String())
	}
}
