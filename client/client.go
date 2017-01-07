package client

import (
	"github.com/drgarcia1986/gonews/providers"
	"github.com/drgarcia1986/gonews/story"

	pb "gopkg.in/cheggaaa/pb.v1"
)

type Client struct {
	provider providers.Provider
	withPB   bool
	pb       *pb.ProgressBar
}

func (c *Client) GetStories(storyType, limit int) ([]*story.Story, error) {
	generator, err := c.provider.GetStories(storyType, limit)
	if err != nil {
		return nil, err
	}

	var stories []*story.Story
	c.startProgressBar(limit)
	for future := range generator {
		request := <-future
		if request.Err != nil {
			return nil, err
		}
		stories = append(stories, request.Story)
		c.incrementProgressBar()
	}
	c.finishProgressBar()
	return stories, nil
}

func (c *Client) startProgressBar(count int) {
	if c.withPB {
		c.pb = pb.StartNew(count)
	}
}

func (c *Client) incrementProgressBar() {
	if c.pb != nil {
		c.pb.Increment()
	}
}

func (c *Client) finishProgressBar() {
	if c.pb != nil {
		c.pb.Finish()
	}
}

func factory(provider providers.Provider, withPB bool) *Client {
	return &Client{provider: provider, withPB: withPB}
}

func New(provider providers.Provider) *Client {
	return factory(provider, false)
}

func NewWithPB(provider providers.Provider) *Client {
	return factory(provider, true)
}
