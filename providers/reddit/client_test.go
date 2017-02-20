package reddit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drgarcia1986/gonews/story"
	"github.com/drgarcia1986/gonews/utils"
)

func TestName(t *testing.T) {
	client := New()
	expected := "reddit"
	if name := client.Name(); name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestSubRedditName(t *testing.T) {
	client := NewSubReddit("golang")
	expected := "reddit-golang"
	if name := client.Name(); name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}
}

func TestGetURL(t *testing.T) {
	var getUrlTests = []struct {
		storyType int
		limit     int
		subReddit string
		expected  string
	}{
		{story.TopStories, 5, "", fmt.Sprintf("%s/top.json?limit=5", urlBase)},
		{story.NewStories, 3, "", fmt.Sprintf("%s/new.json?limit=3", urlBase)},
		{story.TopStories, 10, "golang", fmt.Sprintf("%s/r/golang/top.json?limit=10", urlBase)},
	}

	for _, tt := range getUrlTests {
		actual := getURL(tt.storyType, tt.limit, tt.subReddit)
		if actual != tt.expected {
			t.Errorf(
				"getUrl(%d, %d, %s): expected %s, actual %s",
				tt.storyType, tt.limit, tt.subReddit,
				tt.expected, actual,
			)
		}
	}
}

func TestMakeRequestWithUserAgentHeader(t *testing.T) {
	var userAgent string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.Header.Get("User-Agent")
	}))
	defer ts.Close()

	expectedUserAgent := fmt.Sprintf("gonews:v%s (by /u/drgarcia1986)", utils.Version)
	makeRequest(ts.URL)
	if userAgent != expectedUserAgent {
		t.Errorf("Expected %s, got %s", expectedUserAgent, userAgent)
	}

}

func TestGetStories(t *testing.T) {
	expectedTitle := "test"
	expectedURL := "http://test.com"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{"data": {"children": [
				{"data": {
					"title": "%s",
					"url": "%s"
				}}
			]}}`, expectedTitle, expectedURL)
	}))
	defer ts.Close()

	stories, err := getStories(ts.URL)
	if err != nil {
		t.Errorf("Error on get stories %v", err)
	}

	if len(stories) != 1 {
		t.Errorf("Expected 1, got %d", len(stories))
	}

	story := stories[0]
	if story.Title != expectedTitle {
		t.Errorf("Expected %s, got %s", expectedTitle, story.Title)
	}
	if story.URL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, story.URL)
	}
}

func TestGetStoriesGenerator(t *testing.T) {
	expectedTitles := []string{"test", "test 2"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{"data": {"children": [
				{"data": {
					"title": "%s",
					"url": "http://test.com"
				}},
				{"data": {
					"title": "%s",
					"url": "http://test2.com"
				}}
			]}}`, expectedTitles[0], expectedTitles[1])
	}))
	defer ts.Close()

	urlBase = ts.URL
	client := New()

	generator, err := client.GetStories(story.TopStories, 2)
	if err != nil {
		t.Errorf("Error on get stories %v", err)
	}

	i := 0
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
