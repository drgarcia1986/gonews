package client

import (
	"testing"

	"github.com/drgarcia1986/gonews/progressbar"
	"github.com/drgarcia1986/gonews/providers"
	"github.com/drgarcia1986/gonews/story"
)

func TestGetStories(t *testing.T) {
	expectedStories := []*story.Story{
		{"Foo", "http://golang.org"},
		{"Bar", "http://google.com"},
	}

	fake := &providers.Fake{expectedStories}
	c := New(fake, progressbar.New())

	stories, err := c.GetStories(1, 1)
	if err != nil {
		t.Errorf("error on get stories: %v", err)
	}

	for i, s := range stories {
		if s.Title != expectedStories[i].Title {
			t.Errorf("Expected %s, got %s", expectedStories[i].Title, s.Title)
		}
		if s.URL != expectedStories[i].URL {
			t.Errorf("Expected %s, got %s", expectedStories[i].URL, s.URL)
		}
	}
}
