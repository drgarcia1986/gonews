package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/drgarcia1986/gonews/providers"
	"github.com/drgarcia1986/gonews/story"
	"github.com/drgarcia1986/gonews/utils"
)

const urlBase = "https://www.reddit.com/"

var (
	sufixUrlTopStories = "top.json"
	sufixUrlNewStories = "new.json"
)

type redditResponse struct {
	Data struct {
		Children []struct {
			Data *story.Story `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Reddit struct {
	subReddit string
}

func getURL(storyType, limit int, subReddit string) string {
	url := urlBase
	if subReddit != "" {
		url = fmt.Sprintf("%s/r/%s", url, subReddit)
	}

	switch storyType {
	case story.TopStories:
		url = fmt.Sprintf("%s/%s", url, sufixUrlTopStories)
	default:
		url = fmt.Sprintf("%s/%s", url, sufixUrlNewStories)
	}

	return fmt.Sprintf("%s?limit=%d", url, limit)
}

func makeRequest(url string) (*http.Response, error) {
	client := new(http.Client)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("gonews:v%s (by /u/drgarcia1986)", utils.Version))
	return client.Do(req)
}

func getStories(url string) ([]*story.Story, error) {
	resp, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var rr redditResponse
	if err = json.Unmarshal(body, &rr); err != nil {
		return nil, err
	}

	result := make([]*story.Story, 0)
	for _, s := range rr.Data.Children {
		result = append(result, s.Data)
	}

	return result, nil
}

func (r *Reddit) GetStories(storyType, limit int) (<-chan chan *providers.StoryRequest, error) {
	url := getURL(storyType, limit, r.subReddit)
	stories, err := getStories(url)
	if err != nil {
		return nil, err
	}

	generator := make(chan chan *providers.StoryRequest, len(stories))
	go func() {
		defer close(generator)
		for _, s := range stories {
			f := make(chan *providers.StoryRequest, 1)
			f <- &providers.StoryRequest{s, nil}
			close(f)

			generator <- f
		}
	}()

	return generator, nil
}

func (r *Reddit) Name() string {
	if r.subReddit == "" {
		return "reddit"
	}
	return fmt.Sprintf("reddit-%s", r.subReddit)
}

func NewSubReddit(subReddit string) providers.Provider {
	return &Reddit{subReddit}
}

func New() providers.Provider {
	return new(Reddit)
}
