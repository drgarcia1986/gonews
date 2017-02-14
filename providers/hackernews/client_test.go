package hackernews

import (
	"fmt"
	"testing"

	"github.com/drgarcia1986/gonews/story"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestHackernewsName(t *testing.T) {
	client := New()
	expected := "HackerNews"
	if name := client.Name(); name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestHackernewsGetStoryUrl(t *testing.T) {
	expectedURL := fmt.Sprintf("%s/1.json", urlStoryBase)

	if url := getStoryURL(1); url != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, url)
	}
}

func TestGetStoryIds(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", urlTopStories,
		httpmock.NewStringResponder(200, "[1111, 2222, 3333]"),
	)

	expectedIds := []int{1111, 2222, 3333}
	storyIds, err := getStoryIds(urlTopStories)
	if err != nil {
		t.Errorf("Error on get story ids %v", err)
	}
	for i, value := range storyIds {
		if value != expectedIds[i] {
			t.Errorf("Expected %v, got %v", expectedIds, storyIds)
		}
	}
}

func TestGetStory(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := getStoryURL(1)
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200,
			`{
				"by": "gopher", "descendants": 0, "id": 1,
				"score": 1, "time": 1482804640, "title": "test",
				"type": "story", "url": "http://test.com"
			}`),
	)

	story, err := getStory(1)
	if err != nil {
		t.Errorf("Error on get story ids %v", err)
	}

	expectedTitle := "test"
	expectedURL := "http://test.com"
	if story.Title != expectedTitle {
		t.Errorf("Expected %s, got %s", expectedTitle, story.Title)
	}
	if story.URL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, story.URL)
	}
}

func TestGetUrl(t *testing.T) {
	if url := getURL(story.NewStories); url != urlNewStories {
		t.Errorf("Expected %s, got %s", urlNewStories, url)
	}

	if url := getURL(story.TopStories); url != urlTopStories {
		t.Errorf("Expected %s, got %s", urlTopStories, url)
	}
}

func TestGetStories(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", urlTopStories,
		httpmock.NewStringResponder(200, "[1, 2, 3]"),
	)
	url := getStoryURL(1)
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200,
			`{
				"by": "gopher", "descendants": 0, "id": 1,
				"score": 1, "time": 1482804640, "title": "test",
				"type": "story", "url": "http://test.com"
			}`),
	)
	url = getStoryURL(2)
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200,
			`{
				"by": "gopher", "descendants": 0, "id": 2,
				"score": 1, "time": 1482804640, "title": "test 2",
				"type": "story", "url": "http://test2.com"
			}`),
	)

	client := New()
	generator, err := client.GetStories(story.TopStories, 2)
	if err != nil {
		t.Errorf("Error on get stories ids %v", err)
	}

	i := 0
	expectedTitles := []string{"test", "test 2"}
	for future := range generator {
		r := <-future

		if r.Err != nil {
			t.Errorf("Error on get future stories: %v", r.Err)
		}

		if r.Story.Title != expectedTitles[i] {
			t.Errorf("Expected %s, got %s", expectedTitles[i], r.Story.Title)
		}
		i++
	}
}
