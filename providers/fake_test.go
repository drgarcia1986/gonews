package providers

import (
	"testing"

	"github.com/drgarcia1986/gonews/story"
)

func TestFakeName(t *testing.T) {
	fake := &Fake{[]*story.Story{}}
	expected := "fake"
	if name := fake.Name(); name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestFakeGetStories(t *testing.T) {
	stories := []*story.Story{
		{"Foo", "http://golang.org"},
		{"Bar", "http://google.com"},
	}

	fake := &Fake{stories}
	generator, err := fake.GetStories(0, 0)
	if err != nil {
		t.Errorf("Error on get fake stories: %v", err)
	}

	i := 0
	for future := range generator {
		r := <-future

		if r.Err != nil {
			t.Errorf("Error on get future stories: %v", r.Err)
		}

		expectedTitle := fake.Stories[i].Title
		expectedURL := fake.Stories[i].URL

		if r.Story.Title != expectedTitle {
			t.Errorf("Expected %s, got %s", expectedTitle, r.Story.Title)
		}

		if r.Story.URL != expectedURL {
			t.Errorf("Expected %s, got %s", expectedURL, r.Story.URL)
		}
		i++
	}
}
