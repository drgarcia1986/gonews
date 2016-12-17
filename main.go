package main

import (
	"flag"
	"fmt"

	"github.com/drgarcia1986/gonews/gui"
	"github.com/drgarcia1986/gonews/hackernews"
)

var (
	limit         int
	storyType     int
	flagStoryType string
)

func init() {
	flag.IntVar(&limit, "limit", 10, "Number of Stories to get")
	flag.StringVar(&flagStoryType, "type", "new", "Stories Type ('new' or 'top')")

	flag.Parse()

	if flagStoryType == "new" {
		storyType = hackernews.NewStories
	} else {
		storyType = hackernews.TopStories
	}
}

func main() {
	fmt.Println("Getting HackerNews stories")

	hn := hackernews.New()
	stories, err := hn.GetStories(storyType, limit)
	if err != nil {
		panic(err)
	}

	g := gui.New(stories)
	if err = g.Run(); err != nil {
		panic(err)
	}
}
