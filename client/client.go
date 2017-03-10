package client

import (
	"github.com/drgarcia1986/gonews/progressbar"
	"github.com/drgarcia1986/gonews/providers"
	"github.com/drgarcia1986/gonews/story"
)

type Client struct {
	provider providers.Provider
	pb       progressbar.ProgressBar
}

func (c *Client) GetStories(storyType, limit int) ([]*story.Story, error) {
	generator, err := c.provider.GetStories(storyType, limit)
	if err != nil {
		return nil, err
	}

	var stories []*story.Story
	c.pb.Start(limit)
	for future := range generator {
		request := <-future
		if request.Err != nil {
			return nil, err
		}
		stories = append(stories, request.Story)
		c.pb.Increment()
	}
	c.pb.Finish()
	return stories, nil
}

func New(provider providers.Provider, progressbar progressbar.ProgressBar) *Client {
	return &Client{provider: provider, pb: progressbar}
}
