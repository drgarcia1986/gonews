package providers

import "github.com/drgarcia1986/gonews/story"

type StoryRequest struct {
	Story *story.Story
	Err   error
}

type Provider interface {
	Name() string
	GetStories(int, int) (<-chan chan *StoryRequest, error)
}
