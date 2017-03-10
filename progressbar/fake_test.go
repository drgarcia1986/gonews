package progressbar

import "testing"

func TestFakeStart(t *testing.T) {
	p := new(FakeProgressBar)
	p.Start(1)

	if !p.started {
		t.Error("expected progressbar started")
	}
}

func TestFakeIncrement(t *testing.T) {
	p := new(FakeProgressBar)
	p.Increment()

	if p.incCount != 1 {
		t.Errorf("expected 1, got %d", p.incCount)
	}
}

func TestFakeFinish(t *testing.T) {
	p := new(FakeProgressBar)
	p.Finish()

	if !p.finished {
		t.Error("expected progressbar finished")
	}
}
