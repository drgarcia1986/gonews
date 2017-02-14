package providers

import "github.com/drgarcia1986/gonews/story"

type Fake struct {
	Stories []*story.Story
}

func (f *Fake) GetStories(storyType, limit int) (<-chan chan *StoryRequest, error) {
	generator := make(chan chan *StoryRequest, len(f.Stories))
	go func() {
		defer close(generator)
		for _, s := range f.Stories {
			future := make(chan *StoryRequest, 1)
			future <- &StoryRequest{s, nil}
			close(future)

			generator <- future
		}
	}()
	return generator, nil
}

func (f *Fake) Name() string {
	return "fake"
}

func NewFake() Provider {
	stories := make([]*story.Story, 0)
	return &Fake{stories}
}
