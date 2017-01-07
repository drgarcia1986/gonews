package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/drgarcia1986/gonews/gui"
	"github.com/drgarcia1986/gonews/hackernews"
)

const version = "0.0.1"

var (
	limit         int
	storyType     int
	flagStoryType string
)

func init() {
	flag.IntVar(&limit, "limit", 10, "Number of Stories to get")
	flag.StringVar(&flagStoryType, "type", "top", "Stories Type ('new' or 'top')")

	flag.Parse()

	if len(flag.Args()) > 0 && flag.Args()[0] == "version" {
		fmt.Printf("GoNews %s\n", version)
		os.Exit(0)
	}

	if flagStoryType == "new" {
		storyType = hackernews.NewStories
	} else {
		storyType = hackernews.TopStories
	}
}

func main() {
	fmt.Println("Getting HackerNews stories")

	hn := hackernews.NewWithPB()
	stories, err := hn.GetStories(storyType, limit)
	if err != nil {
		panic(err)
	}

	g := gui.New(stories)
	if err = g.Run(); err != nil {
		panic(err)
	}
}
