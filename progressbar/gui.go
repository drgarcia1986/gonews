package progressbar

import pb "gopkg.in/cheggaaa/pb.v1"

type GUIProgressBar struct {
	p *pb.ProgressBar
}

func (g *GUIProgressBar) Start(count int) {
	g.p = pb.StartNew(count)
}

func (g *GUIProgressBar) Increment() {
	g.p.Increment()
}

func (g *GUIProgressBar) Finish() {
	g.p.Finish()
}

func New() ProgressBar {
	return new(GUIProgressBar)
}
