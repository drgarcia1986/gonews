package hackernews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct{}
type Story struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}
type storyRequest struct {
	id       int
	response *Story
	err      error
}

const (
	urlBase    = "https://hacker-news.firebaseio.com/v0"
	TopStories = iota
	NewStories
)

var (
	urlTopStories = fmt.Sprintf("%s/topstories.json", urlBase)
	urlNewStories = fmt.Sprintf("%s/newstories.json", urlBase)
	urlStoryBase  = fmt.Sprintf("%s/item", urlBase)
)

func getStoryUrl(id int) string {
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

func getStory(id int) (*Story, error) {
	url := getStoryUrl(id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var story Story

	if err = json.Unmarshal(body, &story); err != nil {
		return nil, err
	}
	return &story, nil
}

func getUrl(storyType int) string {
	switch storyType {
	case TopStories:
		return urlTopStories
	default:
		return urlNewStories
	}
}

func storiesGenerator(targetIds []int) chan chan *storyRequest {
	generator := make(chan chan *storyRequest, len(targetIds))

	go func() {
		for _, id := range targetIds {
			generator <- func(id int) chan *storyRequest {
				future := make(chan *storyRequest, 1)
				go func() {
					story, err := getStory(id)
					future <- &storyRequest{id, story, err}
					close(future)
				}()
				return future
			}(id)
		}
		close(generator)
	}()

	return generator
}

func (c *Client) GetStories(storyType, limit int) ([]*Story, error) {
	stories := []*Story{}
	url := getUrl(storyType)
	ids, err := getStoryIds(url)
	if err != nil {
		return nil, err
	}
	targetIds := ids[:limit]

	storiesMap := make(map[int]*Story)
	for future := range storiesGenerator(targetIds) {
		request := <-future
		if request.err != nil {
			return nil, err
		}
		storiesMap[request.id] = request.response
	}

	// To keep order
	for _, id := range targetIds {
		stories = append(stories, storiesMap[id])
	}
	return stories, nil
}

func New() *Client {
	return &Client{}
}
