package hackernews

import (
	"fmt"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestGetStoryUrl(t *testing.T) {
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
	if url := getURL(NewStories); url != urlNewStories {
		t.Errorf("Expected %s, got %s", urlNewStories, url)
	}

	if url := getURL(TopStories); url != urlTopStories {
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
	stories, err := client.GetStories(TopStories, 2)
	if err != nil {
		t.Errorf("Error on get stories ids %v", err)
	}
	if len(stories) != 2 {
		t.Errorf("Expected to get 2 stories, got %d", len(stories))
	}

	if stories[0].Title != "test" {
		t.Errorf("Expected test, got %s", stories[0].Title)
	}

	if stories[1].Title != "test 2" {
		t.Errorf("Expected test 2, got %s", stories[1].Title)
	}
}
