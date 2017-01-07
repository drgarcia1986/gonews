package hackernews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/drgarcia1986/gonews/providers"
	"github.com/drgarcia1986/gonews/story"
)

const urlBase = "https://hacker-news.firebaseio.com/v0"

var (
	urlTopStories = fmt.Sprintf("%s/topstories.json", urlBase)
	urlNewStories = fmt.Sprintf("%s/newstories.json", urlBase)
	urlStoryBase  = fmt.Sprintf("%s/item", urlBase)
)

type HackerNews struct{}

func getStoryURL(id int) string {
	return fmt.Sprintf("%s/%d.json", urlStoryBase, id)
}

func getStoryIds(url string) ([]int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ids := []int{}
	if err = json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func getStory(id int) (*story.Story, error) {
	url := getStoryURL(id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var story story.Story

	if err = json.Unmarshal(body, &story); err != nil {
		return nil, err
	}
	return &story, nil
}

func getURL(storyType int) string {
	switch storyType {
	case story.TopStories:
		return urlTopStories
	default:
		return urlNewStories
	}
}

func storiesGenerator(targetIds []int) <-chan chan *providers.StoryRequest {
	generator := make(chan chan *providers.StoryRequest, len(targetIds))

	go func() {
		for _, id := range targetIds {
			generator <- func(id int) chan *providers.StoryRequest {
				future := make(chan *providers.StoryRequest, 1)
				go func() {
					story, err := getStory(id)
					future <- &providers.StoryRequest{story, err}
					close(future)
				}()
				return future
			}(id)
		}
		close(generator)
	}()

	return generator
}

func (h *HackerNews) GetStories(storyType, limit int) (<-chan chan *providers.StoryRequest, error) {
	url := getURL(storyType)
	ids, err := getStoryIds(url)
	if err != nil {
		return nil, err
	}
	targetIds := ids[:limit]

	return storiesGenerator(targetIds), nil
}

func New() providers.Provider {
	return new(HackerNews)
}
