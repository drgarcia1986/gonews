package progressbar

type FakeProgressBar struct {
	started  bool
	incCount int
	finished bool
}

func (f *FakeProgressBar) Start(count int) {
	f.started = true
}

func (f *FakeProgressBar) Increment() {
	f.incCount++
}

func (f *FakeProgressBar) Finish() {
	f.finished = true
}

func NewFake() ProgressBar {
	return new(FakeProgressBar)
}
