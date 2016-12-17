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
	json.Unmarshal(body, &ids)
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
	json.Unmarshal(body, &story)
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

func (c *Client) GetStories(storyType, limit int) ([]*Story, error) {
	stories := []*Story{}
	url := getUrl(storyType)
	ids, err := getStoryIds(url)
	targetIds := ids[:limit]
	if err != nil {
		return nil, err
	}

	storyResponses := make(chan *storyRequest, limit)
	for _, id := range targetIds {
		go func(id int) {
			story, err := getStory(id)
			storyResponses <- &storyRequest{id, story, err}
		}(id)
	}

	storiesMap := make(map[int]*Story)
	for i := 0; i < limit; i++ {
		request := <-storyResponses
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
