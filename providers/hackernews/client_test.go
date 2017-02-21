package hackernews

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drgarcia1986/gonews/story"
)

func TestHackernewsName(t *testing.T) {
	client := New()
	expected := "HackerNews"
	if name := client.Name(); name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
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

func TestHackernewsGetStoryUrl(t *testing.T) {
	expectedURL := fmt.Sprintf("%s/1.json", urlStoryBase)

	if url := getStoryURL(1); url != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, url)
	}
}

func TestGetStoryIds(t *testing.T) {
	expectedIds := []int{1111, 2222, 3333}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "[%d, %d, %d]", expectedIds[0], expectedIds[1], expectedIds[2])
	}))
	defer ts.Close()

	storyIds, err := getStoryIds(ts.URL)
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
	expectedTitle := "test"
	expectedURL := "http://test.com"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"by": "gopher", "descendants": 0, "id": 1,
				"score": 1, "time": 1482804640, "title": "%s",
				"type": "story", "url": "%s"
			}`, expectedTitle, expectedURL)
	}))
	defer ts.Close()

	urlStoryBase = ts.URL
	story, err := getStory(1)
	if err != nil {
		t.Errorf("Error on get story ids %v", err)
	}

	if story.Title != expectedTitle {
		t.Errorf("Expected %s, got %s", expectedTitle, story.Title)
	}
	if story.URL != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, story.URL)
	}
}

func TestGetStories(t *testing.T) {
	urlTS := "/top.json"
	mux := http.NewServeMux()
	mux.HandleFunc(urlTS, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "[1, 2, 3]")
	})

	expectedTitles := []string{"test", "test 2"}
	mux.HandleFunc("/1.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"by": "gopher", "descendants": 0, "id": 1,
				"score": 1, "time": 1482804640, "title": "%s",
				"type": "story", "url": "http://test.com"
			}
		`, expectedTitles[0])
	})
	mux.HandleFunc("/2.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"by": "gopher", "descendants": 0, "id": 2,
				"score": 1, "time": 1482804640, "title": "%s",
				"type": "story", "url": "http://test.com"
			}
		`, expectedTitles[1])
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	urlStoryBase = ts.URL
	urlTopStories = fmt.Sprintf("%s/%s", ts.URL, urlTS)

	client := New()
	generator, err := client.GetStories(story.TopStories, 2)
	if err != nil {
		t.Errorf("Error on get stories ids %v", err)
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
